package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Notification struct {
	IsOpen        *bool `bson:"isOpen"`
	IsOpenLong    *bool `bson:"isOpenLong"`
	IsExtremeTemp *bool `bson:"isTextreme"`
}

type Reading struct {
	Id            primitive.ObjectID `bson:"_id"`
	TimeStamp     time.Time          `bson:"time"`
	Temp          int                `bson:"temperature"`
	*Notification `bson:"notifications"`
}

func makeReading(books *mongo.Collection, ctx context.Context) {
	minTemp := -30
	maxTemp := 20
	isOpen := false
	isOpenLong := false
	isExtremeTemp := false

	temp := rand.Intn(maxTemp-minTemp) + minTemp
	dOpenT := rand.Intn(30-1) + 1

	if temp < -20 || maxTemp < 15 {
		isExtremeTemp = true
	}

	if dOpenT > 20 {
		isOpenLong = true
	}

	r := Reading{
		Id:        primitive.NewObjectID(),
		TimeStamp: time.Now(),
		Temp:      temp,
		Notification: &Notification{
			IsOpen:        &isOpen,
			IsOpenLong:    &isOpenLong,
			IsExtremeTemp: &isExtremeTemp}}

	_, err := books.InsertOne(ctx, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func prepEnv() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("failed to load env")
		return
	}
}

func handleReadings(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1/readings/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	mongoURI := "mongodb+srv://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASSWORD") + "@cluster0.zkhjg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}

	ctx := context.Background()
	err = client.Connect(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	var reading Reading
	var readings []Reading

	b, err := client.Database("Cluster0").Collection("Readings").Find(ctx, bson.M{})

	for b.Next(ctx) {
		err := b.Decode(&reading)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		readings = append(readings, reading)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(readings)
}

func spawnReadings(wg *sync.WaitGroup) {
	go func() {
		mongoURI := "mongodb+srv://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASSWORD") + "@cluster0.zkhjg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

		client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		err = client.Connect(ctx)

		if err != nil {
			log.Fatal(err)
		}

		defer client.Disconnect(ctx)

		col := client.Database("Cluster0").Collection("Readings")

		i := 0

		for i < 600 {
			time.Sleep(time.Second * 5)
			makeReading(col, ctx)
			fmt.Printf("new reading created: %v\n", i)
			i++
		}

		fmt.Println("Spawning has been stopped for common sense reasons")

		wg.Done()
	}()
}

func server(wg *sync.WaitGroup) {
	func() {
		defer wg.Done()
		http.HandleFunc("/api/v1/readings/", handleReadings)
		log.Fatal(http.ListenAndServe(":5050", nil))
	}()
}

func main() {
	prepEnv()

	wg := new(sync.WaitGroup)

	// Should've passed db connection from here avoiding duplicate code, making it configurable and faster
	// Couldn't find a way to make it work in the reasonable time with parallel processes
	// decided to go forth with the current solution

	// to avoid potential collision
	wg.Add(2)

	go spawnReadings(wg)

	go server(wg)

	wg.Wait()
}
