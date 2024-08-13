package handlers

import (
	"context"
	"net/http"

	"github.com/Panji-Utama/chat-app-backend/models"
	"github.com/Panji-Utama/chat-app-backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var usersCollection *mongo.Collection

func Register(client *mongo.Client) gin.HandlerFunc {
    usersCollection = client.Database("chat_app").Collection("users")

    return func(c *gin.Context) {
        var creds models.Credentials
        if err := c.BindJSON(&creds); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
            return
        }

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            return
        }

        _, err = usersCollection.InsertOne(context.TODO(), bson.M{
            "email": creds.Email,
            "username": creds.Username,
            "password": string(hashedPassword),
        })
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
    }
}

func Login(client *mongo.Client) gin.HandlerFunc {
	usersCollection = client.Database("chat_app").Collection("users")

	return func(c *gin.Context) {
		var creds models.Login
		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		var storedUser models.Credentials
		err := usersCollection.FindOne(context.TODO(), bson.M{"email": creds.Email}).Decode(&storedUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email does not exist"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		token, err := utils.GenerateJWT(storedUser.Id, storedUser.Email, storedUser.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.SetCookie("token", token, 3600, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
	}
}



func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clear the cookie by setting a negative max-age
		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}
}