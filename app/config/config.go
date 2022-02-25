package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func envMongoURI() string {
	// dir, oldErr := os.Getwd()
	// if oldErr != nil {
	// 	log.Fatal(oldErr)
	// }
	// fmt.Println(dir)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv("MONGOURI")
}

func connectDB()  *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(envMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

    //ping the database
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    } else {
		fmt.Println("Connected to MongoDB")
	}

	return client
}

var AlbumCollection = connectDB().Database("greatest-album").Collection("album")
var AccountCollection = connectDB().Database("greatest-album").Collection("account")