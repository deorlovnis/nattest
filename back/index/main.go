package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

type Reading struct {
	Id        primitive.ObjectID  `bson:"_id" json:"_id"`
	TimeStamp primitive.Timestamp `bson:"time"`
	Temp      float32             `bson:"temperature"`
	isOpen    bool                `bson:"isOpen"`
}

func makeReading(books *mongo.Collection, ctx context.Context) {
	// TODO: make data random
	r := Reading{Id: primitive.NewObjectID(), Temp: 36.6, isOpen: false}

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
		log.Fatal(err)
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	var reading Reading
	var readings []Reading

	b, err := client.Database("Cluster0").Collection("Readings").Find(ctx, bson.M{})

	for b.Next(ctx) {
		err := b.Decode(&reading)
		if err != nil {
			return
		}
		readings = append(readings, reading)
	}

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
		// to avoid infinite spawning and limited it to 10 mins if you need to present this project
		for i < 600 {
			time.Sleep(time.Second)
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
	// decided to go forth with the current solution for now

	// to avoid potential collision
	wg.Add(2)

	// go spawnReadings(wg)

	go server(wg)

	wg.Wait()
}
