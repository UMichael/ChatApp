package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Message ...
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

//This is the websocket interface
func handleconnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	clients[ws] = true
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}
func handlemessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}

		}
	}
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go handlemessages()
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	http.HandleFunc("/ws", handleconnection)
	fmt.Println("hello")
	http.ListenAndServe(":"+os.Getenv("HTTP\_PLATFORM\_PORT"), nil)
}
