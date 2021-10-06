package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
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

func makeReading(books *mongo.Collection, Ctx context.Context) {
	r := Reading{Id: primitive.NewObjectID(), Temp: 36.6, isOpen: false}

	_, err := books.InsertOne(Ctx, r)
	if err != nil {
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

func spawnReadings(wg *sync.WaitGroup) {
	go func() {
		mongoURI := "mongodb+srv://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASSWORD") + "@cluster0.zkhjg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

		fmt.Println(os.Getenv("MONGO_USER"))
		client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)

		db := client.Database("Cluster0")
		ReadingsCollection := db.Collection("Readings")

		i := 0
		// to avoid infinite spawning and limited it to 10 mins if you need to present it
		for i < 600 {
			time.Sleep(time.Second)
			makeReading(ReadingsCollection, ctx)
			fmt.Printf("new reading created: %v\n", i)
			i++
		}

		fmt.Println("Spawning has been stopped for common sense reasons")
		wg.Done()
	}()
}

func main() {
	prepEnv()

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go spawnReadings(wg)

	wg.Wait()
}
