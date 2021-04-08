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

func TestRegistrationMessageLoad(t *testing.T) {
	// Start channels broadcasting
	roomStorage = RoomStorage{items: make(map[int]Room)}
	clientsMap = ClientsMap{items: make(map[UserUUID]User)}
	var broadcast = make(chan BroadcastRequest)
	var roomNotificationQueue = make(chan BroadcastRequest, 1000)
	startMessagingRoutines(broadcast, roomNotificationQueue)
	defer closeBroadcastChannels(broadcast, roomNotificationQueue)

	// Start test
	payload := IncomingMessage{
		Username: "pengiun",
		MsgType:  "register",
	}

	volume := 1000
	threads := 5
	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads;i++ {
		go func() {
			for j := 0; j < volume; j++ {
				s, ws := newWSServer(t, handleConnections)
				sendMessage(t, ws, payload)
				msg := receiveWSMessage(t, ws)
				if msg.Type != "room_list" {
					t.Error("Expected room list")
				}
				ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				s.Close()
				ws.Close()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestRoomJoin(t *testing.T) {
	roomStorage = RoomStorage{items: make(map[int]Room)}
	clientsMap = ClientsMap{items: make(map[UserUUID]User)}

	// Three test connections
	s, ws := newWSServer(t, handleConnections)
	s_lobby, ws_lobby := newWSServer(t, handleConnections)
	s_joined, ws_joined := newWSServer(t, handleConnections)

	// Start channels broadcasting
	var broadcast = make(chan BroadcastRequest)
	var roomNotificationQueue = make(chan BroadcastRequest, 1000)
	startMessagingRoutines(broadcast, roomNotificationQueue)
	defer closeBroadcastChannels(broadcast, roomNotificationQueue)

	// Start test
	payload := IncomingMessage{
		Username: "pengiun",
		MsgType:  "register",
	}

	// Register all three users
	sendMessage(t, ws, payload)
	sendMessage(t, ws_lobby, payload)
	sendMessage(t, ws_joined, payload)
	_ = receiveWSMessage(t, ws)
	_ = receiveWSMessage(t, ws_lobby)
	_ = receiveWSMessage(t, ws_joined)

	// Create a room
	payload = IncomingMessage{
		RoomName: "Test",
		MsgType:  "create_room",
		Lat: 42,
		Lng: 42,
	}
	sendMessage(t, ws, payload)
	msg := receiveWSMessage(t, ws_lobby)
	if msg.Type != "room_list" {
		t.Error("Expected room list")
	}

	// Second user joins the room
	payload = IncomingMessage{
		RoomID: 0,
		MsgType:  "join_room",
	}
	sendMessage(t, ws_joined, payload)
	msg = receiveWSMessage(t, ws)
	if msg.Type != "room_update" {
		t.Error("Expected room update")
	}

	// Second user leaves the room
	payload = IncomingMessage{
		MsgType:  "leave_room",
	}
	sendMessage(t, ws_joined, payload)
	msg = receiveWSMessage(t, ws)
	if msg.Type != "room_update" {
		t.Error("Expected room update")
	}

	// Close all sockets
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	s.Close()
	ws.Close()

	ws_lobby.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	s_lobby.Close()
	ws_lobby.Close()

	ws_joined.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	s_joined.Close()
	ws_joined.Close()
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