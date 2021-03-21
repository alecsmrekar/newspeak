package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Room struct {
	ID int
	Name string
	Location GeoLocation
	Members []UserUUID
}

// Concurrency safe
type RoomStorage struct {
	sync.RWMutex
	items map[int]Room
}

// Quasi factory for Rooms
func createRoom (name string, lat coordinate, lng coordinate) Room {
	new := Room{
		Name:     name,
		Location: GeoLocation{lat,lng},
	}
	roomStorage.RegisterRoom(&new)
	return new
}

// Add member to room
func (room *Room) addMember (id UserUUID) {
	*room = roomStorage.AddMember(room.ID, id)
}

// Register new room
func (data *RoomStorage) RegisterRoom(room *Room) {
	data.Lock()
	defer data.Unlock()
	nextID := len(data.items)
	room.ID = nextID
	data.items[nextID] = *room
}

// Get the room object
func (data *RoomStorage) GetRoom(id int) (Room, bool) {
	data.Lock()
	defer data.Unlock()
	room, ok := data.items[id]
	return room, ok
}

// Add a member to a room
func (data *RoomStorage) AddMember(ID int, uid UserUUID) Room {
	data.Lock()
	defer data.Unlock()
	room := data.items[ID]
	room.Members = append(room.Members, uid)
	return room
}

// Remove member from room
func (data *RoomStorage) RemoveMember(ID int, member *websocket.Conn) {
	data.Lock()
	defer data.Unlock()
	room := data.items[ID]
	room.Members = append(room.Members, member)
}

// Delete room
func (data *RoomStorage) Delete(room *Room) {
	data.Lock()
	defer data.Unlock()
	delete(data.items, room.ID)
}

// Get all rooms
func (data *RoomStorage) GetAll() map[int]Room {
	data.Lock()
	defer data.Unlock()
	return (*data).items
}

func (data *RoomStorage) GetListAll() []Room {
	var list []Room
	all := data.GetAll()
	for _, value := range all {
		list = append(list, value)
	}
	return list
}