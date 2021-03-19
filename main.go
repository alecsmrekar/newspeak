package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients_map = ClientsMap{items: make(map[*websocket.Conn]UserData)}
var broadcast = make(chan Broadcast)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Broadcast struct {
	Username string     `json:"username"`
	Message  string     `json:"message"`
}

type radius int
type coordinate float32

type RegisterMessage struct {
	Username string
	Radius radius
	Lat coordinate
	Lng coordinate
}

type Message struct {
	Username string     `json:"username"`
	Message  string     `json:"message"`
	CoordLat coordinate `json:"lat"`
	CoordLng coordinate `json:"lng"`
	MsgType  string     `json:"type"`
	Radius   radius     `json:"radius"`
}

type UserPayload struct {
	message Message
}

func attachClient (clients *ClientsMap, connectionKey *websocket.Conn) {
	clients.Set(connectionKey, UserData{})
}

func detachClient (clients *ClientsMap, connectionKey *websocket.Conn) {
	clients.Delete(connectionKey)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer ws.Close()
	attachClient(&clients_map, ws)
	user := getUserByConnectionID(ws)

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			detachClient(&clients_map, ws)
			break
		}

		switch msg.MsgType {
		case "message":
			broadcast <- Broadcast{
				Message:  msg.Message,
				Username: user.data.Username,
			}
		case "register":
			user.updateData(UserPayload{message: msg}, &register{})
		case "radius":
			user.updateData(UserPayload{message: msg}, &updateRadius{})
		case "location":
			user.updateData(UserPayload{message: msg}, &updateCoordinates{})
		default:
			log.Println("Unknown communication type")
			detachClient(&clients_map, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		// grab next message from the broadcast channel
		msg := <-broadcast
		// send it out to every client that is currently connected
		for KeyValPair := range clients_map.Iter() {
			client := KeyValPair.Key
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				detachClient(&clients_map, client)
			}
		}
	}
}

func main() {

	fs := http.FileServer(http.Dir("./web/dist"))
	http.Handle("/", fs)

	// the function will launch a new goroutine for each request
	http.HandleFunc("/ws", handleConnections)
	for i := 0; i < 1; i++ {
		go handleMessages()
	}

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
