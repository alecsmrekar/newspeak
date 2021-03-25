package main

import (
	"github.com/gorilla/websocket"
	"testing"
)

func TestInitUserData(t *testing.T) {
	var con *websocket.Conn
	user := User{
		connectionKey: con,
	}
	initUserData(&user, "test2")
	if user.username != "test2" || user.currentRoom != -1 {
		t.Error("Failed user initialization")
	}
}


func TestRemovingUserFromRoomStorage(t *testing.T) {
	var con *websocket.Conn
	roomStorage = RoomStorage{items: make(map[int]Room)}
	clientsMap = ClientsMap{items: make(map[UserUUID]User)}
	user := User{
		connectionKey: con,
		currentRoom: 1,
	}
	var uuid UserUUID = user.connectionKey
	clientsMap.Set(uuid, user)

	room := Room{
		ID:         1,
		Members:    []UserUUID{uuid},
	}
	roomStorage.RegisterRoom(&room)
	leaveRoom(uuid)
	room, _ = roomStorage.GetRoom(1)
	if len(room.Members) > 0 {
		t.Error("Error removing user from clients map")
	}
}

