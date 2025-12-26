package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connection() *mongo.Client {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Unable to load env file.")
	}

	MongoDB := os.Getenv("MONGODB_URI")
	if MongoDB == "" {
		log.Fatal("Unable to find the MONGODB_URI")
	}
	fmt.Println("MONGODB_URI: ", MongoDB)

	clientOptions := options.Client().ApplyURI(MongoDB)
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil
	}
	return client
}

var Client *mongo.Client = Connection()

func OpenCollection(collectionName string) *mongo.Collection {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Unable to load the env file")
	}

	databaseName := os.Getenv("DATABASE_NAME")

	fmt.Println("DATABASE_NAME: ", databaseName)

	collection := Client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		return nil
	}

	return collection
}
