package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type album struct {
	Number int
	Year int
	Album string
	Artist string
	Genre string
	Subgenre string
}

func connectDB()  *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27017"))
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

var collection = connectDB().Database("greatest-album").Collection("album")

func main() {
	r := gin.Default()

	connectDB()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message" : "hello world!",
		})
	})

	r.GET("/albums", getAlbums())

	r.GET("/album/:number", getAlbum())

	r.Run()
}

func getAlbums() gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor, err := collection.Find(context.TODO(), bson.D{})
		if err != nil {
			log.Fatal(err)
		}

		var albums []album

		for cursor.Next(context.TODO()) {
			var album album
			if err = cursor.Decode(&album); err != nil {
				log.Fatal(err)
			}

			albums = append(albums, album)
		}

		c.JSON(200, albums)
	}
}

func getAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		number := c.Param("number")
		newNumber, _ := strconv.Atoi(number)

		var result album
		err := collection.FindOne(context.TODO(), bson.M{"Number": newNumber}).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(200, result)
	}
}