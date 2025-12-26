package main

import (
	"fmt"

	"github.com/agarwal-parth/Movie-Streaming-App/Server/MagicStreamMoviesServer/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	// fmt.Println("Hello, MagicStreamMoviesServer!")
	router := gin.Default()

	router.GET("hello", func(c *gin.Context) {
		c.String(200, "Hello, MagicStreamMoviesServer!")
	})

	router.GET("movies", controllers.GetMovies())
	router.GET("movie/:imdb_id", controllers.GetMovie())
	router.POST("addmovie", controllers.AddMovie())
	router.POST("register", controllers.RegisterUser())

	err := router.Run(":8000")
	if err != nil {
		fmt.Println("Failed to start Server")
	}
}
