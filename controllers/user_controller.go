package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/agarwal-parth/Movie-Streaming-App/Server/MagicStreamMoviesServer/database"
	"github.com/agarwal-parth/Movie-Streaming-App/Server/MagicStreamMoviesServer/models"
	"github.com/agarwal-parth/Movie-Streaming-App/Server/MagicStreamMoviesServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection("users")

func HashPassword(password string) (string, error) {
	hashedPassWord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassWord), err
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.ShouldBindBodyWithJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input Data!"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation of Input Failed.", "details": err.Error()})
			return
		}

		hashedPassWord, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash password", "details": err.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User with the same email exists"})
			return
		}

		user.UserID = bson.NewObjectID().Hex()
		user.Password = hashedPassWord
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add user in the Server"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		var userLogin models.UserLogin
		if err := c.ShouldBindBodyWithJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input Data"})
			return
		}

		var foundUser models.User
		if err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&foundUser); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password!!"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password!"})
			return
		}

		token, RefreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate tokens."})
			return
		}

		err = utils.UpdateAllTokens(foundUser.UserID, token, RefreshToken)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update the tokens."})
			return
		}

		c.JSON(http.StatusOK, models.UserResponse{
			UserID:          foundUser.UserID,
			FirstName:       foundUser.FirstName,
			LastName:        foundUser.LastName,
			Email:           foundUser.Email,
			Role:            foundUser.Role,
			Token:           token,
			RefreshToken:    RefreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		})
	}
}
