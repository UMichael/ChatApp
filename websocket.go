package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	EnableCompression: true,
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
		//name, p, errs := ws.ReadMessage()
		//fmt.Println(name, p, errs)
		fmt.Println("connection ", broadcast)
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
		fmt.Println("handle ", msg)
	}
}
func main() {
	//this is how we server directory
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	//like google.com/ws you get going to run it now
	http.HandleFunc("/ws", handleconnection)
	go handlemessages()
	http.ListenAndServe(":8000", nil)
}
