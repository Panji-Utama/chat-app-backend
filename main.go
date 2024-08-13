package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Panji-Utama/chat-app-backend/router"
	"github.com/Panji-Utama/chat-app-backend/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Get MongoDB URI from environment variable
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        log.Fatalf("MONGODB_URI environment variable is not set")
    }

    // Set up MongoDB connection
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    // Initialize router
    r := router.InitializeRouter(client)

    // Start WebSocket handler for messages
    go utils.HandleMessages()

    // Run the server
    log.Println("Server started at :8000")
    if err := r.Run(":7000"); err != nil {
        log.Fatalf("Server Run Failed: %v", err)
    }
}
