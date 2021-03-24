package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
)

type BroadcastRequest struct {
	broadcast OutgoingBroadcast
	receivers []*websocket.Conn
}

// Represents an outgoing chat message
type OutgoingBroadcast struct {
	Type string				`json:"type"`
	Username string     	`json:"username"`
	Message  string     	`json:"message"`
	RoomName string			`json:"room_name"`
	RoomList map[int]Room	`json:"room_list"`
	RoomID int 				`json:"room_id"`
	RoomUsers []string		`json:"users"`
}

// The way we identify the user is modular
// Currently, it's based on the WebSocker connection ID
type UserUUID *websocket.Conn
type coordinate float32

// Represents an incoming communication
type IncomingMessage struct {
	RoomID 	int				`json:"room_id"`
	RoomName 	string		`json:"room_name"`
	Username string     	`json:"username"`
	Message  string     	`json:"message"`
	MsgType  string     	`json:"type"`
	Lat 	coordinate 		`json:"lat"`
	Lng 	coordinate 		`json:"lng"`
}

type GeoLocation struct {
	Lat coordinate
	Lng coordinate
}

var roomNotificationQueue = make(chan BroadcastRequest, 1000)
var roomStorage = RoomStorage{items: make(map[int]Room)}
var clientsMap = ClientsMap{items: make(map[UserUUID]User)}
var lobby = ConcurrentSlice{}
var broadcast = make(chan BroadcastRequest)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func attachClient (clients *ClientsMap, connectionKey UserUUID) {
	clients.Set(connectionKey, User{
		connectionKey: connectionKey,
	})
}

func detachClient (clients *ClientsMap, id UserUUID) {
	leaveRoom(id)
	clients.Delete(id)
}

// A thread that handles one client's communications
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer ws.Close()
	attachClient(&clientsMap, ws)

	var uuid UserUUID = ws
	_, found := clientsMap.Get(uuid)
	if !found {
		log.Println("Error adding client to client map")
		ws.Close()
		detachClient(&clientsMap, uuid)
	}

	// Test, create a sample room
	for i := 0; i < 5; i++ {
		nr := rand.Intn(50000)
		c1 := coordinate(rand.Intn(80))
		c2 := coordinate(rand.Intn(80))

		sampleRoom := Room{
			Name:     fmt.Sprintf("Room %v", nr),
			Location: GeoLocation{c1, c2},
		}
		roomStorage.RegisterRoom(&sampleRoom)
	}


	for {
		var msg IncomingMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			leaveRoom(uuid)
			log.Printf("error: %v", err)
			detachClient(&clientsMap, ws)
			break
		}

		manager := CommunicationsManager{msg, uuid}

		switch msg.MsgType {
		case "message":
			// Send msg to room
			manager.sendMsg()
		case "create_room":
			manager.createRoom()
		case "join_room":
			manager.joinRoomProcess(msg.RoomID)
		case "register":
			manager.register()
		case "leave_room":
			manager.leaveRoom()
		default:
			log.Println("Unknown communication type")
			detachClient(&clientsMap, ws)
			break
		}
	}
}

// Receives messages from the broadcast channel and sends them out
func handleMessageBroadcasting() {
	for {
		// grab next message from the broadcast channel
		request := <-broadcast
		for _, user := range request.receivers {
			sendBroadcast(user, request.broadcast)
		}

	}
}

// Receives newly created rooms and and notifies clients
func handleRoomNotifications() {
	for {
		request := <-roomNotificationQueue
		for _, user := range request.receivers {
			sendBroadcast(user, request.broadcast)
		}
	}
}

// Send a message to a single user
func sendBroadcast(client *websocket.Conn, msg OutgoingBroadcast) {
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		client.Close()
		detachClient(&clientsMap, client)
	}
}

func main() {

	fs := http.FileServer(http.Dir("./web/dist"))
	http.Handle("/", fs)

	// the function will launch a new goroutine for each request
	http.HandleFunc("/ws", handleConnections)

	// Launch a few thread that send out messages
	for i := 0; i < 10; i++ {
		go handleMessageBroadcasting()
	}

	// Launch a few thread that send out room notifications
	for i := 0; i < 10; i++ {
		go handleRoomNotifications()
	}

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
