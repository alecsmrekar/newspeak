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
	MembersFull []User
	MembersCount int
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

// Get a lightweight version of the object
func (room *Room) getRoomProxy () Room {
	return Room{
		ID:           room.ID,
		Name:         room.Name,
		Location:     room.Location,
		MembersCount: len(room.Members),
	}
}

// Get a lightweight version of the object
func (room *Room) getRoomWithMembers () Room {
	var members []User
	for _, uuid  := range room.Members {
		user, ok := clients_map.Get(uuid)
		if ok {
			members = append(members, user)
		}
	}
	return Room{
		ID:           room.ID,
		Name:         room.Name,
		Location:     room.Location,
		MembersFull: members,
	}
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

// Get all rooms
func (data *RoomStorage) GetAllProxied() map[int]Room {
	rooms := data.GetAll()
	proxies := make(map[int]Room)
	for id, room := range rooms {
		proxy := Room{
			ID:           room.ID,
			Name:         room.Name,
			Location:     room.Location,
			MembersCount: len(room.Members),
		}
		proxies[id] = proxy
	}
	return proxies
}

func (data *RoomStorage) GetListAll() []Room {
	var list []Room
	all := data.GetAll()
	for _, value := range all {
		list = append(list, value)
	}
	return list
}