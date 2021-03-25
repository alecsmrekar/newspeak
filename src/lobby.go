package main

import (
	"sync"
)

type ConcurrentSlice struct {
	sync.RWMutex
	items []UserUUID
}

// Appends an item to the concurrent slice
func (cs *ConcurrentSlice) Set(uuid UserUUID) {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, uuid)
}

// Remove a user from the lobby
func (cs *ConcurrentSlice) Delete(uuid UserUUID) {
	cs.Lock()
	defer cs.Unlock()
	for index, id := range cs.items {
		if id == uuid {
			cs.items[index] = cs.items[len(cs.items)-1]
			cs.items[len(cs.items)-1] = nil
			cs.items = cs.items[:len(cs.items)-1]
			break
		}
	}
}

// Returns all ids
func (cs *ConcurrentSlice) GetAll() []UserUUID {
	cs.Lock()
	defer cs.Unlock()
	return cs.items
}
