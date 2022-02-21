package main

import (
	"greatest-album/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message" : "hello world!",
		})
	})

	r.GET("/albums", controller.GetAlbums())

	r.GET("/album/:number", controller.GetAlbum())

	r.POST("album/", controller.AddAlbum())

	r.DELETE("album/:number", controller.DeleteAlbum())

	r.Run()
}