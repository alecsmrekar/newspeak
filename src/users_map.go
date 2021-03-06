package main

import (
	"sync"
)

// Concurrency safe map of clients
type ClientsMap struct {
	sync.RWMutex
	items map[UserUUID]User
}

// Sets a key in the concurrent map of clients
func (clients *ClientsMap) Set(id UserUUID, value User) {
	clients.Lock()
	defer clients.Unlock()
	clients.items[id] = value
}

// Deletes a key in the concurrent map of clients
func (clients *ClientsMap) Delete (id UserUUID) {
	clients.Lock()
	defer clients.Unlock()
	delete(clients.items, id)
}

// Gets a key from the concurrent map  of clients
func (clients *ClientsMap) Get(id UserUUID) (User, bool) {
	clients.Lock()
	defer clients.Unlock()
	value, ok := clients.items[id]
	return value, ok
}

// Returns a list of User objects based on user UUIDs
func (clients *ClientsMap) LookupUserIDs(ids []UserUUID) []User {
	clients.Lock()
	defer clients.Unlock()
	var users []User
	for _, id := range ids {
		value, ok := clients.items[id]
		if ok {
			users = append(users, value)
		}
	}
	return users
}

// Gets a key from the concurrent map  of clients
func (clients *ClientsMap) AddUserToGroup(id UserUUID, room int) User {
	clients.Lock()
	defer clients.Unlock()
	value, ok := clients.items[id]
	if ok {
		value.currentRoom = room
		clients.items[id] = value
	}
	return value
}

// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (clients *ClientsMap) Iter() <-chan User {

	c := make(chan User)

	f := func() {
		clients.Lock()
		defer clients.Unlock()

		for _, v := range clients.items {
			c <- v
		}
		close(c)
	}
	go f()
	return c
}