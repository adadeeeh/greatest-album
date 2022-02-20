package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type album struct {
	Number int `json:"Number"` 
	Year int `json:"Year" validate:"required"`
	Album string `json:"Album" validate:"required"`
	Artist string `json:"Artist" validate:"required"`
	Genre string `json:"Genre" validate:"required"`
	Subgenre string `json:"Subgenre" validate:"required"`
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

	r.POST("album/", addAlbum())

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

func addAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		var album album
		
		// bind data to struct
		if err := c.BindJSON(&album); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validation empty field
		if validationErr := validator.New().Struct(&album); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}
		
		opts := options.Update().SetUpsert(true)
		filter := bson.M{"Number": album.Number}
		update := bson.M{"$set": bson.M{	// Using bson to match the capitalize field
			"Number": album.Number,
			"Year": album.Year,
			"Album": album.Album,
			"Artist": album.Artist,
			"Genre": album.Genre,
			"Subgenre": album.Subgenre,
		}}

		result, err := collection.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		if result.MatchedCount != 0 {
			c.JSON(http.StatusCreated, "Matched and replaced an existing document")
		}
		if result.UpsertedCount != 0 {
			c.JSON(http.StatusCreated, result)
		}

		// Using bson to match the capitalize field
		// result, err := collection.InsertOne(context.TODO(), bson.M{
		// 	"Number": album.Number,
		// 	"Year": album.Year,
		// 	"Album": album.Album,
		// 	"Artist": album.Artist,
		// 	"Genre": album.Genre,
		// 	"Subgenre": album.Subgenre,
		// })
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// }

		// c.JSON(http.StatusCreated, result)
	}
}