package controller

import (
	"context"
	"fmt"
	"greatest-album/config"
	"greatest-album/model"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAlbums() gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor, err := config.Collection.Find(context.TODO(), bson.D{})
		if err != nil {
			log.Println(err)
			return
		}

		var results []model.Album

		for cursor.Next(context.TODO()) {
			var result model.Album
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

func GetAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		number := c.Param("number")
		newNumber, _ := strconv.Atoi(number)

		var result model.Album
		err := config.Collection.FindOne(context.TODO(), bson.M{"Number": newNumber}).Decode(&result)
		if err != nil {
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": result,
		})
	}
}

func AddAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		var album model.Album
		
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

		result, err := config.Collection.UpdateOne(context.TODO(), filter, update, opts)
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

func DeleteAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		number := c.Param("number")
		newNumber, _ := strconv.Atoi(number)
		filter := bson.M{"Number": newNumber}

		result, err := config.Collection.DeleteOne(context.TODO(), filter)
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