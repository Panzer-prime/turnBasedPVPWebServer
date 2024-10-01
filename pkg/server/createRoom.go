package server

import (
	"arena-game/pkg/game"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	RoomID     string
	RoomName   string
	Password   string
	Clients    map[*websocket.Conn]bool
	Broadcast  chan []byte
	MaxPlayers int
	GameState  *game.GameState
}

type RoomManager struct {
	Rooms map[string]*Room
	mu    sync.RWMutex
}

var roomManager = RoomManager{
	Rooms: make(map[string]*Room),
}

func (rm *RoomManager) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method; POST required", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var req struct {
		RoomName     string `json:"roomName"`
		RoomPassword string `json:"roomPassword"`
		UserName     string `json:"userName"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(req.RoomPassword)
	if err != nil {
		http.Error(w, "Internal server error: password hashing failed", http.StatusInternalServerError)
		return
	}

	roomID := CreateUniqueID()

	newRoom := &Room{
		RoomID:     roomID,
		RoomName:   req.RoomName,
		Password:   hashedPassword,
		Clients:    make(map[*websocket.Conn]bool),
		MaxPlayers: 2,
		GameState: &game.GameState{
			Players: []*game.Player{
				game.CreatePlayer(CreateUniqueID(), req.UserName),
			},
		},
	}

	rm.mu.Lock()
	rm.Rooms[roomID] = newRoom
	rm.mu.Unlock()

	//go HandleMessages(conn, room)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"roomID": roomID,
		"user": map[string]string{
			"id":   newRoom.GameState.Players[0].ID,
			"name": newRoom.GameState.Players[0].Name,
		},
		"message": "room was created",
	})
}
