package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBGetClient() (*mongo.Client, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error: cannot load .env file!")
	}
	mongo_uri := os.Getenv("MONGODB_URL")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_uri))
	if err != nil {
		panic(err)
	}

	return client, err
}

func DBDisconnect(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func AddNewDocumentForTest(client *mongo.Client, user interface{}) {
	coll := client.Database("user-tokens").Collection("JWT")
	if _, err := coll.InsertOne(context.TODO(), user); err != nil {
		panic(err)
	}
}