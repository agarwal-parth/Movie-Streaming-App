package routes

import (
	"github.com/agarwal-parth/Movie-Streaming-App/Server/MagicStreamMoviesServer/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine) {
	router.GET("movies", controllers.GetMovies())
	router.POST("register", controllers.RegisterUser())
	router.POST("/login", controllers.LoginUser())
}
