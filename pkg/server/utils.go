package server

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUniqueID() string {
	return uuid.New().String()
}

func hashPassword(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (rm *RoomManager) GetRoomById(roomID string) (*Room, error) {
	rm.mu.RLock()
	room, exists := rm.Rooms[roomID]
	rm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("Room %s does not exist", roomID)
	}

	return room, nil
}

