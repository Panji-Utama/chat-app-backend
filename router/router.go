package router

import (
	"github.com/gin-gonic/gin"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()

	// Auth routes
	router.POST("/api/login", handlers.Login)
	router.POST("/api/register", handlers.Register)

	// Websocket
	router.GET("/ws", handlers.HandleConnections)

	return router
}