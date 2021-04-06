package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestWSConnectionV2(t *testing.T) {
	// Start WS Socker Handler
	s, ws := newWSServer(t, handleConnections)
	defer s.Close()
	defer ws.Close()

	// Start channels broadcasting
	var wg sync.WaitGroup
	wg.Add(2)
	go handleMessageBroadcasting(&wg)
	go handleRoomNotifications(&wg)
	defer closeBroadcastChannels()
	wg.Wait()

	// Start test
	payload := IncomingMessage{
		Username: "pengiun",
		MsgType:  "register",
	}

	sendMessage(t, ws, payload)
	msg := receiveWSMessage(t, ws)
	if msg.Type != "room_list" {
		t.Error("Extected room list")
	}
}

func newWSServer(t *testing.T, h http.HandlerFunc) (*httptest.Server, *websocket.Conn) {
	t.Helper()

	s := httptest.NewServer(h)
	wsURL := httpToWS(t, s.URL)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	return s, ws
}

func httpToWS(t *testing.T, u string) string {
	t.Helper()

	wsURL, err := url.Parse(u)
	if err != nil {
		t.Fatal(err)
	}

	switch wsURL.Scheme {
	case "http":
		wsURL.Scheme = "ws"
	case "https":
		wsURL.Scheme = "wss"
	}

	return wsURL.String()
}

func sendMessage(t *testing.T, ws *websocket.Conn, msg IncomingMessage) {
	t.Helper()

	m, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.WriteMessage(websocket.BinaryMessage, m); err != nil {
		t.Fatalf("%v", err)
	}
}

func receiveWSMessage(t *testing.T, ws *websocket.Conn) OutgoingBroadcast {
	t.Helper()

	_, m, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	var reply OutgoingBroadcast
	err = json.Unmarshal(m, &reply)
	if err != nil {
		t.Fatal(err)
	}

	return reply
}