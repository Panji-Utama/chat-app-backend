package utils

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Panji-Utama/chat-app-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Declaring the variables
var (
    clients   = make(map[*websocket.Conn]string) // Connected clients
    broadcast = make(chan models.Message)        // Broadcast channel
    upgrader  = websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
)

var messagesCollection *mongo.Collection

// SetClient initializes the MongoDB client in this package
func SetClient(mongoClient *mongo.Client) {
    messagesCollection = mongoClient.Database("chat_app").Collection("messages")
}

func HandleConnections(c *gin.Context) {
    ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Fatalf("Failed to set websocket upgrade: %v", err)
    }
    defer ws.Close()

    clients[ws] = "" // Add the connection to the list of clients

    for {
        var msg models.Message
        err := ws.ReadJSON(&msg)
        if err != nil {
            log.Printf("Error reading message: %v", err)
            delete(clients, ws)
            break
        }
        msg.Timestamp = primitive.NewDateTimeFromTime(time.Now())

        // Save message to MongoDB
        _, err = messagesCollection.InsertOne(context.TODO(), msg)
        if err != nil {
            log.Printf("Failed to save message: %v", err)
            continue
        }

        // Send the message to the broadcast channel
        broadcast <- msg
    }
}

func HandleMessages() {
    for {
        msg := <-broadcast
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Printf("Error writing message: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}
