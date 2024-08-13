package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/Panji-Utama/chat-app-backend/utils"
)

func HandleConnections(c *gin.Context) {
    utils.HandleConnections(c)
}