package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Room struct {
	ID int
	Name string
	Location GeoLocation
	Members []*websocket.Conn
}

func createRoom (name string, lat coordinate, lng coordinate) Room {
	new := Room{
		Name:     name,
		Location: GeoLocation{lat,lng},
	}
	roomStorage.RegisterRoom(&new)
	return new
}

func (room *Room) addMember (conn *websocket.Conn) {
	roomStorage.AddMember(room.ID, conn)
}

// Concurrency safe
type RoomStorage struct {
	sync.RWMutex
	items map[int]Room
}

// Sets a key in the concurrent map of clients
func (data *RoomStorage) RegisterRoom(room *Room) {
	data.Lock()
	defer data.Unlock()
	nextID := len(data.items)
	room.ID = nextID
	data.items[nextID] = *room
}

// Sets a key in the concurrent map of clients
func (data *RoomStorage) GetRoom(id int) (Room, bool) {
	data.Lock()
	defer data.Unlock()
	room, ok := data.items[id]
	return room, ok
}

func (data *RoomStorage) AddMember(ID int, member *websocket.Conn) {
	data.Lock()
	defer data.Unlock()
	room := data.items[ID]
	room.Members = append(room.Members, member)
}

func (data *RoomStorage) Delete(room *Room) {
	data.Lock()
	defer data.Unlock()
	delete(data.items, room.ID)
}

func (data *RoomStorage) GetAll() map[int]Room {
	data.Lock()
	defer data.Unlock()
	return (*data).items
}