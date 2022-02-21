package main

import (
	"greatest-album/route"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message" : "hello world!",
		})
	})

	route.AlbumRoute(r)

	r.Run()
}