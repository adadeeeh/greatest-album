package main

import (
	"context"
	"greatest-album/config"
	"greatest-album/model"
	"greatest-album/route"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	r := gin.Default()

	account := getAccount()

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		account["username"].(string): account["password"].(string),
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message" : "hello world!",
		})
	})

	route.AlbumRoute(authorized)

	r.Run()
}

func getAccount() map[string]interface{} {
	var result model.Account
	err := config.AccountCollection.FindOne(context.TODO(), bson.D{}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	account := map[string]interface{} {
		"username": result.Username,
		"password": result.Password,
	}

	return account

}