package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
)


// Represents an outgoing chat message
type OutgoingBroadcast struct {
	Type string			`json:"type"`
	Username string     `json:"username"`
	Message  string     `json:"message"`
	Room Room			`json:"room"`
	RoomList map[int]Room		`json:"room_list"`
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
var lobby = ClientsMap{items: make(map[UserUUID]User)} // This can be converted into an array
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

func detachClient (clients *ClientsMap, id UserUUID) {
	leaveRoom(id)
	//TODO notify room members
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
	attachClient(&clients_map, ws)

	var uuid UserUUID = ws
	var usersCurrentRoomID = -1
	user, found := clients_map.Get(uuid)
	if !found {
		log.Println("Error adding client to client map")
		ws.Close()
		detachClient(&clients_map, uuid)
	}

	// Test, create a sample room
	for i := 0; i < 20; i++ {
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
			detachClient(&clients_map, ws)
			break
		}

		switch msg.MsgType {
		case "message":
			// Send msg to room
			broadcast <- OutgoingBroadcast{
				Type: "message",
				Message:  msg.Message,
				Username: user.username,
			}
		case "create_room":
			// Create room
			// Remove user from lobby
			// Add him to the room
			createdRoom := createRoom(msg.RoomName, msg.Lat, msg.Lng)
			createdRoom.addMember(uuid)
			lobby.Delete(uuid)
			usersCurrentRoomID = createdRoom.ID
			reply := OutgoingBroadcast{
				Type:     "joined_room",
				Room:     createdRoom.getRoomWithMembers(),
			}
			sendBroadcast(ws, reply)
			roomNotificationQueue <- createdRoom
		case "join_room":
			// Remove user from lobby
			// Add him to the room
			// to all room members: send them list of members
			if usersCurrentRoomID >= 0 {
				roomStorage.RemoveMember(usersCurrentRoomID, uuid)
			} else {
				lobby.Delete(uuid)
			}
			joinedRoom := roomStorage.AddMember(msg.RoomID, uuid)
			reply := OutgoingBroadcast{
				Type:     "room_status",
				Room:     joinedRoom.getRoomWithMembers(),
			}
			sendBroadcast(ws, reply)
		case "register":
			// Add him to table of users
			// Send him a list of all rooms
			strategy := &register{}
			strategy.update(&user, UserPayload{message: msg})
			clients_map.Set(uuid, user)
			lobby.Set(uuid, user)
			reply := OutgoingBroadcast{
				Type:     "room_list",
				RoomList:     roomStorage.GetAllProxied(),
			}
			sendBroadcast(ws, reply)
		case "leave_room":
			// Leave room
			// Add to lobby
			// Send him a list of all rooms
			leaveRoom(uuid)
			lobby.Set(uuid, user)
			reply := OutgoingBroadcast{
				Type:     "room_list",
				RoomList:     roomStorage.GetAllProxied(),
			}
			sendBroadcast(ws, reply)
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
