// handlers/chat.go
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Panji-Utama/chat-app-backend/models"
	"github.com/Panji-Utama/chat-app-backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleConnections(c *gin.Context) {
    utils.HandleConnections(c)
}

func GetMessages(client *mongo.Client) gin.HandlerFunc {
    messagesCollection := client.Database("chat_app").Collection("messages")

    return func(c *gin.Context) {
        sender := c.Query("sender")
        recipient := c.Query("recipient")

        var messages []models.Message
        filter := bson.M{
            "$or": []bson.M{
                {"sender": sender, "recipient": recipient},
                {"sender": recipient, "recipient": sender},
            },
        }
        cur, err := messagesCollection.Find(context.TODO(), filter)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
            return
        }

        if err := cur.All(context.TODO(), &messages); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode messages"})
            return
        }

        c.JSON(http.StatusOK, messages)
    }
}

func SaveMessage(client *mongo.Client) gin.HandlerFunc {
    messagesCollection := client.Database("chat_app").Collection("messages")

    return func(c *gin.Context) {
        var msg models.Message
        if err := c.BindJSON(&msg); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
            return
        }

        // Set the timestamp for the message
        msg.Timestamp = primitive.NewDateTimeFromTime(time.Now())

        // Save the message to the MongoDB collection
        _, err := messagesCollection.InsertOne(context.TODO(), msg)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Message saved successfully"})
    }
}
