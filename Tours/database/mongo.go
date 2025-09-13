package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var TourCollection *mongo.Collection
var ReviewCollection *mongo.Collection
var KeyPointCollection *mongo.Collection
var PositionCollection *mongo.Collection

func Init() {
	uri := os.Getenv("MONGO_URI")

	// if uri == "" {
	// 	uri = "mongodb://localhost:27017"
	// }

	dbName := os.Getenv("MONGO_DB")

	// if dbName == "" {
	// 	dbName = "Tours"
	// }

	collName := os.Getenv("MONGO_COLLECTION")
	// if collName == "" {
	// 	collName = "Tours"
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Cannot connect to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	fmt.Println("MongoDB connected")

	Client = client
	TourCollection = client.Database(dbName).Collection(collName)
	ReviewCollection = client.Database(dbName).Collection("Reviews")
	KeyPointCollection = client.Database(dbName).Collection("Keypoints")
	PositionCollection = client.Database(dbName).Collection("Positions")
}
