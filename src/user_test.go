package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
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

func TestWSConnection(t *testing.T) {
	go startWebServer()
	var wg sync.WaitGroup
	wg.Add(2)
	go handleMessageBroadcasting(&wg)
	go handleRoomNotifications(&wg)
	wg.Wait()
	fmt.Println("Test setup ready")

	u := "ws://localhost:8000/ws"
	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	payload := IncomingMessage{
		Username: "pengiun",
		MsgType:  "register",
	}

	if err := ws.WriteJSON(payload); err != nil {
		t.Fatalf("%v", err)
	}

	var msg IncomingMessage
	err = ws.ReadJSON(&msg)
	if err != nil || msg.MsgType != "room_list" {
		t.Fatalf("%v", err)
	}
}
