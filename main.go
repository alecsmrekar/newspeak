package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type UserData struct {
	Lat float32
	Lng float32
	Radius int
}

// Concurrency safe map of clients
type ClientsMap struct {
	sync.RWMutex
	items map[*websocket.Conn]UserData
}

// Concurrent map of clients
type ClientsMapItem struct {
	Key   *websocket.Conn
	Value UserData
}

// Sets a key in the concurrent map of clients
func (clients *ClientsMap) Set(connectionKey *websocket.Conn, value UserData) {
	clients.Lock()
	defer clients.Unlock()
	clients.items[connectionKey] = value
}

// Sets a key in the concurrent map of clients
func (clients *ClientsMap) Delete (connectionKey *websocket.Conn) {
	clients.Lock()
	defer clients.Unlock()
	delete(clients.items, connectionKey)
}

// Gets a key from the concurrent map  of clients
func (clients *ClientsMap) Get(connectionKey *websocket.Conn) (UserData, bool) {
	clients.Lock()
	defer clients.Unlock()
	value, ok := clients.items[connectionKey]
	return value, ok
}


// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (clients *ClientsMap) Iter() <-chan ClientsMapItem {

	c := make(chan ClientsMapItem)

	f := func() {
		clients.Lock()
		defer clients.Unlock()

		for k, v := range clients.items {
			c <- ClientsMapItem{k, v}
		}
		close(c)
	}
	go f()
	return c
}


var clients_map = ClientsMap{items: make(map[*websocket.Conn]UserData)}
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Username 	string `json:"username"`
	Message  	string `json:"message"`
	CoordLat  	float32 `json:"lat"`
	CoordLng  	float32 `json:"lng"`
	Radius  	int `json:"radius"`
	IsRadiusUpdate bool `json:"radiusUpdate"`
	IsLocationUpdate bool `json:"locationUpdate"`
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

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			detachClient(&clients_map, ws)
			break
		}
		// send the new message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// grab next message from the broadcast channel
		msg := <-broadcast
		// send it out to every client that is currently connected
		for KeyValPair := range clients_map.Iter() {
			client := KeyValPair.Key

			if msg.IsRadiusUpdate {
				updateUserRadius(client, msg.Radius)
			}

			if msg.IsLocationUpdate {
				updateUserCoordinates(client, msg.CoordLat, msg.CoordLng)
			}

			if len(msg.Message) > 0 {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					detachClient(&clients_map, client)
				}
			}
		}
	}
}

func updateUserRadius(connectionKey *websocket.Conn, radius int) {
	user_data, found := clients_map.Get(connectionKey)
	if !found {
		log.Println("Unknown client: Tried to update radius")
		connectionKey.Close()
		detachClient(&clients_map, connectionKey)
	}
	user_data.Radius = radius
	clients_map.Set(connectionKey, user_data)
	fmt.Printf("Setting radius\n")
}

func updateUserCoordinates(connectionKey *websocket.Conn, lat float32, lng float32) {
	user_data, found := clients_map.Get(connectionKey)
	if !found {
		log.Println("Unknown client: Tried to update coordinates")
		connectionKey.Close()
		detachClient(&clients_map, connectionKey)
	}
	user_data.Lat = lat
	user_data.Lng = lng
	clients_map.Set(connectionKey, user_data)
	fmt.Printf("Setting coordinates\n")
}

func main() {

	fs := http.FileServer(http.Dir("./web/dist"))
	http.Handle("/", fs)

	// the function will launch a new goroutine for each request
	http.HandleFunc("/ws", handleConnections)
	for i := 0; i < 2; i++ {
		go handleMessages()
	}

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
