package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

func main() {
	prepEnv()

	mongoURI := "mongodb+srv://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASSWORD") + "@cluster0.zkhjg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

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
	makeReading(ReadingsCollection, ctx)

}
