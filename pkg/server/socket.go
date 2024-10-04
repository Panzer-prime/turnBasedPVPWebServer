package server

import (
	"arena-game/pkg/game"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (rm *RoomManager) handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	defer conn.Close()


	roomID := r.URL.Query().Get("roomID")
	playerID := r.URL.Query().Get("playerID")

	if roomID == "" || playerID == "" {
		log.Println("Missing roomID or playerID in query parameters")
		conn.WriteJSON(map[string]interface{}{
			"message": "roomID and playerID are required",
			"error":   true,
		})
		return
	}
	
	

	if err := rm.joinRoom(roomID, conn, playerID); err != nil {
		log.Printf("Could not join room: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	//send game state at the begging of the game
	rm.handleMessages(conn, roomID, playerID)
}

func (rm *RoomManager) joinRoom(roomID string, conn *websocket.Conn, playerID string) error {
	rm.mu.RLock()
	room, err := rm.GetRoomById(roomID)
	rm.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("Room %s not found", roomID)
	}

	// Check if the player is reconnecting
	var playerExists bool
	for _, player := range room.GameState.Players {
		if player.ID == playerID {
			playerExists = true
			break
		}
	}

	if playerExists {
		// Player is reconnecting, allow them to rejoin even if the room is full
		log.Printf("Player %s is reconnecting to room %s", playerID, roomID)

		// If player is already connected, do not add again
		if _, alreadyConnected := room.Clients[conn]; alreadyConnected {
			return fmt.Errorf("player %s is already connected to room %s", playerID, roomID)
		}
	} else {
		// Check if the room is full before adding a new player
		if len(room.GameState.Players) >= room.MaxPlayers {
			conn.WriteJSON(map[string]interface{}{
				"message": "Room is full",
				"error":   true,
			})
			return fmt.Errorf("Room %s is full, max number of players is %d", roomID, room.MaxPlayers)
		}
	}

	// Add or reconnect the player
	rm.mu.Lock()
	room.Clients[conn] = true
	if !playerExists {
		room.GameState.Players = append(room.GameState.Players, game.CreatePlayer(CreateUniqueID(), playerID))
	}
	rm.mu.Unlock()

	log.Printf("Player %s joined/reconnected to room %s", playerID, roomID)
	return nil
}


func (rm *RoomManager) handleMessages(conn *websocket.Conn, roomID string, playerID string)  error {
	rm.mu.RLock()
	room, err := rm.GetRoomById(roomID)
	rm.mu.RUnlock()

	if err != nil {
		log.Printf("Room not found: %v", err)
		return err
	}
	initGameState, err := game.GameInit(room.GameState)

	if err != nil{
		return err
	}
	conn.WriteJSON(initGameState)

	for {

		if game.IsGameOver(room.GameState) {

			conn.WriteJSON(map[string]interface{}{
				"message":   "Game Over",
				"gameState": room.GameState,
			})
			return nil
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client: %v", err)
			return err
		}

		log.Printf("Message received from client: %s and the room is: %s", string(msg), roomID)

		rm.mu.RLock()
		player := room.GameState.Players[room.GameState.TurnIndex]
		rm.mu.RUnlock()

		if playerID != player.ID {
			conn.WriteJSON(map[string]interface{}{
				"message": "Not your turn",
				"error":   true,
			})
			continue
		}

		// Process the player's action (logic to be added here)
		game.HandleAction(msg, room.GameState)
		// Advance the turn after processing
		rm.mu.Lock()
		game.AdvanceTurn(room.GameState)
		rm.mu.Unlock()

		// Broadcast the updated game state to all players
		rm.broadcast(room)
	}
}

func (rm *RoomManager) broadcast(room *Room) {
	for conn := range room.Clients {
		player := room.GameState.Players[room.GameState.TurnIndex]
		err := conn.WriteJSON(map[string]interface{}{
			"message":           "Next turn",
			"currentTurnPlayer": player.ID,
			"gameSate":          room.GameState,
		})
		if err != nil {
			log.Printf("Failed to send turn update to client: %v", err)
		}
	}
}

