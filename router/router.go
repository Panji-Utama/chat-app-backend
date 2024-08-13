package router

import (
	"github.com/Panji-Utama/chat-app-backend/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeRouter(client *mongo.Client) *gin.Engine {
    router := gin.Default()

    router.Use(cors.Default())

    // Group routes under /api prefix
    api := router.Group("/api")
    {
        api.POST("/login", handlers.Login(client))
        api.POST("/register", handlers.Register(client))
        api.POST("/logout", handlers.Logout())  
    }
    
    // WebSocket route, outside /api group if needed
    router.GET("/ws", handlers.HandleConnections)

    return router
}
