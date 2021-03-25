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
