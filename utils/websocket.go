package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
    Sender    string `json:"sender"`
    Recipient string `json:"recipient"`
    Message   string `json:"message"`
}

var clients = make(map[*websocket.Conn]string) // Connected clients and their usernames
var broadcast = make(chan Message)             // Broadcast channel

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func HandleConnections(c *gin.Context) {
    ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Fatalf("Failed to set websocket upgrade: %v", err)
    }
    defer ws.Close()

    // Read initial message to set the username
    var initialMessage Message
    err = ws.ReadJSON(&initialMessage)
    if err != nil {
        log.Printf("Error during initial message read: %v", err)
        return
    }
    clients[ws] = initialMessage.Sender

    for {
        var msg Message
        err := ws.ReadJSON(&msg)
        if err != nil {
            log.Printf("Error reading message: %v", err)
            delete(clients, ws)
            break
        }
        broadcast <- msg
    }
}

func HandleMessages() {
    for {
        msg := <-broadcast
        for client, username := range clients {
            if username == msg.Recipient || username == msg.Sender {
                err := client.WriteJSON(msg)
                if err != nil {
                    log.Printf("Error writing message: %v", err)
                    client.Close()
                    delete(clients, client)
                }
            }
        }
    }
}
