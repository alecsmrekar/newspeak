package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// Quasi User factory
func getUserByConnectionID(connectionKey *websocket.Conn) User {
	user_data, found := clients_map.Get(connectionKey)
	if !found {
		log.Println("Unknown client: Tried to update coordinates")
		connectionKey.Close()
		detachClient(&clients_map, connectionKey)
	}
	return User{
		data:          user_data,
		connectionKey: connectionKey,
	}
}

type User struct {
	data UserData
	connectionKey *websocket.Conn
}

func (user *User) updateData (payload UserPayload, strategy userUpdater) {
	strategy.update(&user.data, payload)
	clients_map.Set(user.connectionKey, user.data)
}

// User Update Interface - Strategy Pattern
type userUpdater interface {
	update(user *UserData, data UserPayload)
}


// User updater - Register User
type register struct {
}

func (l *register) update(user *UserData, data UserPayload) {
	user.Username = data.message.Username
	user.Lat = data.message.CoordLat
	user.Lng = data.message.CoordLng
	user.Radius = data.message.Radius
}

// User updater - Update radius
type updateRadius struct {
}

func (l *updateRadius) update(user *UserData, data UserPayload) {
	user.Radius = data.message.Radius
}

// User updater - Update coordinates
type updateCoordinates struct {
}

func (l *updateCoordinates) update(user *UserData, data UserPayload) {
	user.Lat = data.message.CoordLat
	user.Lng = data.message.CoordLng
}