package main

import "github.com/gin-gonic/gin"

func main() {
	m := gin.Default()
	m.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message" : "hello world!",
		})
	})
	m.Run()
}