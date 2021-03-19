package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type UserData struct {
	Lat coordinate
	Lng coordinate
	Radius radius
	Username string
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