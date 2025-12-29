package main

import (
	"fmt"

	"github.com/agarwal-parth/Movie-Streaming-App/Server/MagicStreamMoviesServer/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// fmt.Println("Hello, MagicStreamMoviesServer!")
	router := gin.Default()

	router.GET("hello", func(c *gin.Context) {
		c.String(200, "Hello, MagicStreamMoviesServer!")
	})
	routes.SetupUnProtectedRoutes(router)
	routes.SetupProtectedRoutes(router)

	err := router.Run(":8000")
	if err != nil {
		fmt.Println("Failed to start Server")
	}
}
