package main

import (
    "log"
    "github.com/Panji-Utama/chat-app-backend/router"
    "github.com/Panji-Utama/chat-app-backend/utils"
)

func main() {
    r := router.InitializeRouter()

    go utils.HandleMessages()

    log.Println("Server started at :8000")
    if err := r.Run(":8000"); err != nil {
        log.Fatalf("Server Run Failed: %v", err)
    }
}
