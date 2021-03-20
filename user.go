package main

import (
	"github.com/gorilla/websocket"
)

type User struct {
	username string
	connectionKey *websocket.Conn
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
}
