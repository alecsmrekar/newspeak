package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

type BroadcastRequest struct {
	broadcast OutgoingBroadcast
	receivers []*websocket.Conn
}

// Represents an outgoing chat message
type OutgoingBroadcast struct {
	Type      string       `json:"type"`
	Username  string       `json:"username"`
	Message   string       `json:"message"`
	RoomName  string       `json:"room_name"`
	RoomList  map[int]Room `json:"room_list"`
	RoomID    int          `json:"room_id"`
	RoomUsers []string     `json:"users"`
}

// The way we identify the user is modular
// Currently, it's based on the WebSocker connection ID
type UserUUID *websocket.Conn
type coordinate float32

// Represents an incoming communication
type IncomingMessage struct {
	RoomID   int        `json:"room_id"`
	RoomName string     `json:"room_name"`
	Username string     `json:"username"`
	Message  string     `json:"message"`
	MsgType  string     `json:"type"`
	Lat      coordinate `json:"lat"`
	Lng      coordinate `json:"lng"`
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

func attachClient(clients *ClientsMap, connectionKey UserUUID) {
	clients.Set(connectionKey, User{
		connectionKey: connectionKey,
	})
}

func detachClient(clients *ClientsMap, uuid UserUUID) {
	leaveRoom(uuid)
	clients.Delete(uuid)
}

// A thread that handles one client's communications
func handleConnections(w http.ResponseWriter, r *http.Request) {
	connectionKey, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer connectionKey.Close()
	attachClient(&clientsMap, connectionKey)

	var uuid UserUUID = connectionKey
	_, found := clientsMap.Get(uuid)
	if !found {
		log.Println("Error adding client to client map")
		connectionKey.Close()
		detachClient(&clientsMap, uuid)
	}

	// ****** Test data: create some sample rooms
	for i := 0; i < 0; i++ {
		nr := rand.Intn(50000)
		c1 := coordinate(rand.Intn(80))
		c2 := coordinate(rand.Intn(80))

		sampleRoom := Room{
			Name:     fmt.Sprintf("Room %v", nr),
			Location: GeoLocation{c1, c2},
		}
		roomStorage.RegisterRoom(&sampleRoom)
	}
	// *****

	for {
		var msg IncomingMessage
		// Read in a new message as JSON and map it to a Message object
		err := connectionKey.ReadJSON(&msg)
		if err != nil {
			leaveRoom(uuid)
			log.Printf("error: %v", err)
			detachClient(&clientsMap, connectionKey)
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
			detachClient(&clientsMap, connectionKey)
			break
		}
	}
}

// Receives messages from the broadcast channel and sends them out
func handleMessageBroadcasting(wg *sync.WaitGroup, chn chan BroadcastRequest) {
	dispatchBroadcast(broadcast, wg)
}

// Receives newly created rooms and and notifies clients
func handleRoomNotifications(wg *sync.WaitGroup, chn chan BroadcastRequest) {
	dispatchBroadcast(roomNotificationQueue, wg)
}

// Dispatches broadcasts from a selected channel
func dispatchBroadcast (chn chan BroadcastRequest, wg *sync.WaitGroup) {
	wg.Done()
	for {
		request, ok := <-chn
		if !ok {
			return
		}
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

func startWebServer() {
	fs := http.FileServer(http.Dir("./web/dist"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)

	// the function will launch a new goroutine for each request
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func startMessagingRoutines(msg chan BroadcastRequest, room chan BroadcastRequest) {
	// Launch a few threads that send out messages
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go handleMessageBroadcasting(&wg, msg)
	}

	// Launch a few threads that send out room notifications
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go handleRoomNotifications(&wg, room)
	}
	wg.Wait()
}

func main() {
	startMessagingRoutines(broadcast, roomNotificationQueue)
	startWebServer()
}

// Closes the broadcast channels and stops the goroutines
func closeBroadcastChannels(msg chan BroadcastRequest, room chan BroadcastRequest) {
	close(msg)
	close(room)
}