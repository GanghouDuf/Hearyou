package ws

import "sync"

type RoomManager struct {
	mu    sync.Mutex
	rooms map[string]*Hub
}

func NewRoomManager() *RoomManager {
	return &RoomManager{rooms: make(map[string]*Hub)}
}

func (room *RoomManager) GetOrCreate(roomName string) *Hub {
	room.mu.Lock()
	defer room.mu.Unlock()
	hub, exists := room.rooms[roomName]

	if exists {
		return hub
	}

	hub = NewHub()
	go hub.Run()
	room.rooms[roomName] = hub
	return hub
}
