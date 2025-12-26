package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/agarwal-parth/Movie-Streaming-App/Server/MagicStreamMoviesServer/database"
	"github.com/agarwal-parth/Movie-Streaming-App/Server/MagicStreamMoviesServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var validate = validator.New()

// If the first letter is capiatalized, it is public other it is private and cannot be exported.
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to fetch the movies from the database"})
		}

		defer cursor.Close(ctx)

		err = cursor.All(ctx, &movies)
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to decode the movies from the database"})
		}

		c.JSON(200, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required."})
			return
		}

		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie Not Found"})
			return
		}

		c.JSON(200, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindBodyWithJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input! Please give input properly."})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input! Not able to validate the input."})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to insert into system.", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}
