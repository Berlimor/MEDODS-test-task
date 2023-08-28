package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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
	db := client.Database("user-tokens")
	allNames, err := db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	exists := false
	for _, name := range allNames {
		if name == "JWT" {
			exists = true
			break
		}
	}
	if !exists{
		db.CreateCollection(context.TODO(), "JWT")
	}
	coll := db.Collection("JWT")
	if _, err := coll.InsertOne(context.TODO(), user); err != nil {
		panic(err)
	}
}