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
var BlogCollection *mongo.Collection
var CommentCollection *mongo.Collection
var LikeCollection *mongo.Collection

func Init() {
	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB")

	if dbName == "" {
		dbName = "Blogs"
	}

	collName := os.Getenv("MONGO_COLLECTION")
	if collName == "" {
		collName = "Blogs"
	}

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

	fmt.Println("MongoDB connected for Blogs")

	Client = client
	BlogCollection = client.Database(dbName).Collection(collName)
	CommentCollection = client.Database(dbName).Collection("Comments")
	LikeCollection = client.Database(dbName).Collection("Likes")
}
