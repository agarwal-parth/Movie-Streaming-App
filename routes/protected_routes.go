package routes

import (
	"github.com/agarwal-parth/Movie-Streaming-App/Server/MagicStreamMoviesServer/controllers"
	"github.com/agarwal-parth/Movie-Streaming-App/Server/MagicStreamMoviesServer/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleWare())
	router.GET("movie/:imdb_id", controllers.GetMovie())
	router.POST("addmovie", controllers.AddMovie())
}
