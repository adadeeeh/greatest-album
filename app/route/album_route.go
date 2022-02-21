package route

import (
	"greatest-album/controller"

	"github.com/gin-gonic/gin"
)

func AlbumRoute(r *gin.Engine) {
	r.GET("/albums", controller.GetAlbums())

	r.GET("/album/:number", controller.GetAlbum())

	r.POST("album/", controller.AddAlbum())

	r.DELETE("album/:number", controller.DeleteAlbum())
}