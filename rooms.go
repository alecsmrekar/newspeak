package main

import (
	"fmt"
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

// Get a list of websocket connections of all room members
func (data *RoomStorage) GetRoomMemberConnections(id int) ([]*websocket.Conn) {
	data.Lock()
	defer data.Unlock()
	room, ok := data.items[id]
	var connections []*websocket.Conn
	if ok {
		for _, member := range room.Members {
			connections = append(connections, member)
		}
	}
	return connections
}

// Add a member to a room
func (data *RoomStorage) AddMember(ID int, uid UserUUID) Room {
	data.Lock()
	defer data.Unlock()
	room := data.items[ID]
	room.Members = append(room.Members, uid)
	data.items[ID] = room
	return room
}

// Remove member from room
// If room is empty, also delete entire room
func (data *RoomStorage) RemoveMember(ID int, member *websocket.Conn) {
	data.Lock()
	defer data.Unlock()
	room := data.items[ID]
	for index, user := range room.Members {
		if user == member {
			room.Members[index] = room.Members[len(room.Members)-1]
			room.Members[len(room.Members)-1] = nil
			room.Members = room.Members[:len(room.Members)-1]
			data.items[ID] = room
			if len(data.items[ID].Members) == 0 {
				data.Delete(&room)
			}
			break
		}
	}
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

// Get a list of all rooms, but a lightweight version for sending to the frontend
func (data *RoomStorage) GetAllProxied() map[int]Room {
	rooms := data.GetAll()
	proxies := make(map[int]Room)
	for id, room := range rooms {
		proxy := Room{
			ID:           room.ID,
			Name:         fmt.Sprintf("%s (%v online)", room.Name,len(room.Members)),
			Location:     room.Location,
		}
		proxies[id] = proxy
	}
	return proxies
}

// Get a list of room member User objects
func getRoomMemberObjects( id int) []User {
	room, ok := roomStorage.GetRoom(id)
	var users []User
	if ok {
		users = clientsMap.LookupUserIDs(room.Members)
	}
	return users
}

// Get a list of room member usernames
func getRoomMemberNames (id int) []string {
	objects := getRoomMemberObjects(id)
	var names []string
	for _, user := range objects {
		names = append(names, user.username)
	}
	return names
}
