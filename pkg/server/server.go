package server

import (
	"fmt"
	"log"
	"net/http"
)

func Server() {
	port := ":3000"
	fmt.Println("server is starting on port: ", port)
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("game")
	})
	http.HandleFunc("/create-room", roomManager.CreateRoom)
	http.HandleFunc("/join-room", roomManager.JoinRoom)
	http.HandleFunc("/connect", roomManager.handleWebSocketConnection)
	fmt.Println()
	log.Fatal(http.ListenAndServe(port, nil))
}
