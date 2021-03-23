package main

import (
	"github.com/gorilla/websocket"
)

type User struct {
	username string
	connectionKey *websocket.Conn
	currentRoom int
}



// User Update Interface - Strategy Pattern
type userUpdater interface {
	update(user *User, data UserPayload)
}

// User updater - Register User
type register struct {
}

func (l *register) update(user *User, data UserPayload) {
	user.username = data.message.Username
	user.currentRoom = -1
}

func leaveRoom(id UserUUID) {
	user, ok := clients_map.Get(id)
	if !ok {
		return
	}
	if user.currentRoom > -1 {
		roomStorage.RemoveMember(user.currentRoom, user.connectionKey)
	}
}