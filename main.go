package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)


// Represents an outgoing chat message
type OutgoingBroadcast struct {
	Type string			`json:"type"`
	Username string     `json:"username"`
	Message  string     `json:"message"`
	Room Room			`json:"room"`
}

// The way we identify the user is modular
// Currently, it's based on the WebSocker connection ID
type UserUUID *websocket.Conn
type coordinate float32

// Represents an incoming communication
type IncomingMessage struct {
	RoomID 	int			`json:"room_id"`
	RoomName 	string		`json:"room_name"`
	Username string     `json:"username"`
	Message  string     `json:"message"`
	MsgType  string     `json:"type"`
	Lat 	coordinate `json:"lat"`
	Lng 	coordinate `json:"lng"`
}

type GeoLocation struct {
	Lat coordinate
	Lng coordinate
}

var roomNotificationQueue = make(chan Room, 1000)
var roomStorage = RoomStorage{items: make(map[int]Room)}
var clients_map = ClientsMap{items: make(map[UserUUID]User)}
var broadcast = make(chan OutgoingBroadcast)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Class that is used to handle user data
type UserPayload struct {
	message IncomingMessage
}

func attachClient (clients *ClientsMap, connectionKey UserUUID) {
	clients.Set(connectionKey, User{
		connectionKey: connectionKey,
	})
}

func detachClient (clients *ClientsMap, connectionKey UserUUID) {
	clients.Delete(connectionKey)
}

// A thread that handles one client's communications
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer ws.Close()
	attachClient(&clients_map, ws)

	var uuid UserUUID = ws
	var usersCurrentRoomID int = -1
	user, found := clients_map.Get(uuid)
	if !found {
		log.Println("Error adding client to client map")
		ws.Close()
		detachClient(&clients_map, uuid)
	}

	for {
		var msg IncomingMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			detachClient(&clients_map, ws)
			break
		}

		switch msg.MsgType {
		case "message":
			broadcast <- OutgoingBroadcast{
				Type: "message",
				Message:  msg.Message,
				Username: user.username,
			}
		case "create_room":
			createdRoom := createRoom(msg.RoomName, msg.Lat, msg.Lng)
			createdRoom.addMember(uuid)
			usersCurrentRoomID = createdRoom.ID
			reply := OutgoingBroadcast{
				Type:     "room_joined",
				Room:     createdRoom,
			}
			sendBroadcast(ws, reply)
			roomNotificationQueue <- createdRoom
		case "join_room":
			if usersCurrentRoomID >= 0 {
				roomStorage.RemoveMember(usersCurrentRoomID, uuid)
			}
			joinedRoom := roomStorage.AddMember(msg.RoomID, uuid)
			reply := OutgoingBroadcast{
				Type:     "room_joined",
				Room:     joinedRoom,
			}
			sendBroadcast(ws, reply)
		case "register":
			strategy := &register{}
			strategy.update(&user, UserPayload{message: msg})
			clients_map.Set(uuid, user)
		default:
			log.Println("Unknown communication type")
			detachClient(&clients_map, ws)
			break
		}
	}
}

// Receives messages from the broadcast channel and sends them out
func handleMessageBroadcasting() {
	for {
		// grab next message from the broadcast channel
		msg := <-broadcast
		// send it out to every client that is currently connected
		for user := range clients_map.Iter() {
			sendBroadcast(user.connectionKey, msg)
		}
	}
}

// Receives newly created rooms and and notifies clients
func handleRoomNotifications() {
	for {
		// grab next room from the queue
		room := <-roomNotificationQueue
		msg := OutgoingBroadcast{
			Type:     "new_room",
			Room:     room,
		}
		// Notify all connected clients
		for user := range clients_map.Iter() {
			sendBroadcast(user.connectionKey, msg)
		}
	}
}

// Send a message to a single user
func sendBroadcast(client *websocket.Conn, msg OutgoingBroadcast) {
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		client.Close()
		detachClient(&clients_map, client)
	}
}

func main() {

	fs := http.FileServer(http.Dir("./web/dist"))
	http.Handle("/", fs)

	// the function will launch a new goroutine for each request
	http.HandleFunc("/ws", handleConnections)

	// Launch a few thread that send out messages
	for i := 0; i < 4; i++ {
		go handleMessageBroadcasting()
	}

	// Launch a few thread that send out room notifications
	for i := 0; i < 4; i++ {
		go handleRoomNotifications()
	}

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
