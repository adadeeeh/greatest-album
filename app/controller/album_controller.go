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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAlbums() gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor, err := config.AlbumCollection.Find(context.TODO(), bson.D{})
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
			"status": http.StatusOK,
			"message": "success",
			"data": results,
		})
	}
}

func GetAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		number := c.Param("number")
		newNumber, _ := strconv.Atoi(number)

		var result model.Album
		err := config.AlbumCollection.FindOne(context.TODO(), bson.M{"Number": newNumber}).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusOK,
					"message": "error",
					"data": fmt.Sprintf("Document with number %v is not found", newNumber),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"message": "error",
				"data": err.Error(),
			})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "success",
			"data": result,
		})
	}
}

func AddAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		var album model.Album
		
		// bind data to struct
		if err := c.BindJSON(&album); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusInternalServerError,
				"message": "error",
				"data": err.Error(),
			})
			log.Println(err)
			return
		}

		// validation empty field
		if validationErr := validator.New().Struct(&album); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"message": "error",
				"data": validationErr.Error(),
			})
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

		result, err := config.AlbumCollection.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"message": "error",
				"data": err.Error(),
			})
			log.Println(err)
			return
		}

		if result.MatchedCount != 0 {
			c.JSON(http.StatusCreated, gin.H{
				"status": http.StatusOK,
				"message": "success",
				"data": "Matched and replaced an existing document",
			})
		}
		if result.UpsertedCount != 0 {
			c.JSON(http.StatusCreated, gin.H{
				"status": http.StatusOK,
				"message": "success",
				"data": "Document added",
			})
		}
	}
}

func DeleteAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		number := c.Param("number")
		newNumber, _ := strconv.Atoi(number)
		filter := bson.M{"Number": newNumber}

		result, err := config.AlbumCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"message": "error",
				"data": err.Error(),
			})
			log.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"message": "success",
				"data": "No document deleted",
			})
		}
		if result.DeletedCount == 1 {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"message": "success",
				"data": fmt.Sprintf("Document with number %v is deleted", newNumber),
			})
		}
	}
}