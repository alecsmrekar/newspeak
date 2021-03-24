package main

import (
	"github.com/gorilla/websocket"
)

type User struct {
	username      string
	connectionKey *websocket.Conn
	currentRoom   int
}

func initUserData(user *User, username string) {
	user.username = username
	user.currentRoom = -1
}

// Removes a user from his current room
func leaveRoom(uuid UserUUID) {
	user, ok := clientsMap.Get(uuid)
	if !ok {
		return
	}
	if user.currentRoom > -1 {
		roomStorage.RemoveMember(user.currentRoom, user.connectionKey)
	}
}
