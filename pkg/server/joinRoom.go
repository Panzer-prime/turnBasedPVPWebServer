package server

import (
	"arena-game/pkg/game"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (rm *RoomManager) JoinRoom(w http.ResponseWriter, r *http.Request) {
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

	var room *Room
	rm.mu.RLock()
	for _, r := range rm.Rooms {
		if r.RoomName == req.RoomName {
			room = r
			break
		}
	}
	rm.mu.RUnlock()

	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(room.Password), []byte(req.RoomPassword)); err != nil {
		http.Error(w, "Invalid room password", http.StatusUnauthorized)
		return
	}

	newPlayer := game.CreatePlayer(CreateUniqueID(), req.UserName)

	rm.mu.Lock()
	defer rm.mu.Unlock()
	room.GameState.Players = append(room.GameState.Players, newPlayer)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"roomID":      room.RoomID,
		"playersList": room.GameState.Players,
	})
}
