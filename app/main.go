package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type album struct {
	Number int `json:"number"` 
	Year int `json:"year" validate:"required"`
	Album string `json:"album" validate:"required"`
	Artist string `json:"artist" validate:"required"`
	Genre string `json:"genre" validate:"required"`
	Subgenre string `json:"subgenre" validate:"required"`
}

func envMongoURI() string {
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

	r.DELETE("album/:number", deleteAlbum())

	r.Run()
}

func getAlbums() gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor, err := collection.Find(context.TODO(), bson.D{})
		if err != nil {
			log.Println(err)
			return
		}

		var results []album

		for cursor.Next(context.TODO()) {
			var result album
			if err = cursor.Decode(&result); err != nil {
				log.Println(err)
				return
			}

			results = append(results, result)
		}

		c.JSON(http.StatusOK, gin.H{
			"data": results,
		})
	}
}

func getAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		number := c.Param("number")
		newNumber, _ := strconv.Atoi(number)

		var result album
		err := collection.FindOne(context.TODO(), bson.M{"Number": newNumber}).Decode(&result)
		if err != nil {
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": result,
		})
	}
}

func addAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		var album album
		
		// bind data to struct
		if err := c.BindJSON(&album); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(err)
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
			log.Println(err)
		}

		if result.MatchedCount != 0 {
			c.JSON(http.StatusCreated, gin.H{"message": "Matched and replaced an existing document."})
		}
		if result.UpsertedCount != 0 {
			c.JSON(http.StatusCreated, gin.H{"message": "Document added."})
		}
	}
}

func deleteAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		number := c.Param("number")
		newNumber, _ := strconv.Atoi(number)
		filter := bson.M{"Number": newNumber}

		result, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Println(err)
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusOK, gin.H{
				"message": "No document deleted.",
			})
		}
		if result.DeletedCount == 1 {
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Document with number %v is deleted.", newNumber),
			})
		}
	}
}