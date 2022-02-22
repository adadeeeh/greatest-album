package route

import (
	"greatest-album/controller"

	"github.com/gin-gonic/gin"
)

func AlbumRoute(authorized *gin.RouterGroup) {
	authorized.GET("/albums", controller.GetAlbums())

	authorized.GET("/album/:number", controller.GetAlbum())

	authorized.POST("album/", controller.AddAlbum())

	authorized.DELETE("album/:number", controller.DeleteAlbum())
}