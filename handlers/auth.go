package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "github.com/Panji-Utama/chat-app-backend/models"
    "github.com/Panji-Utama/chat-app-backend/utils"
)

var users = map[string]string{} // In-memory user storage for simplicity

func Register(c *gin.Context) {
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

    users[creds.Username] = string(hashedPassword)
    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
    var creds models.Credentials
    if err := c.BindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    storedPassword, ok := users[creds.Username]
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := utils.GenerateJWT(creds.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    c.SetCookie("token", token, 3600, "/", "localhost", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
