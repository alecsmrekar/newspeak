package main

import "fmt"

type CommunicationsManager struct {
	IncomingMessage IncomingMessage
	FromUser        UserUUID
}

func (info *CommunicationsManager) sendMsg() {
	user, ok := clientsMap.Get(info.FromUser)
	if ok {
		broadcast <- BroadcastRequest{
			broadcast: OutgoingBroadcast{
				Type:     "message",
				Message:  info.IncomingMessage.Message,
				Username: user.username,
				RoomID:   user.currentRoom,
			},
			receivers: roomStorage.GetRoomMemberConnections(user.currentRoom),
		}
	}
}

func (info *CommunicationsManager) createRoom() {

	// Init room and register in the storage
	createdRoom := createRoom(info.IncomingMessage.RoomName, info.IncomingMessage.Lat, info.IncomingMessage.Lng)

	// Add user to room
	info.joinRoomProcess(createdRoom.ID)

	// Notify lobby of new room
	roomNotificationQueue <- BroadcastRequest{
		broadcast: OutgoingBroadcast{
			Type:     "room_list",
			RoomList: roomStorage.GetAllProxied(),
		},
		receivers: getLobbyUsersConnections(),
	}
}

// Handle the incoming request to join a room
func (info *CommunicationsManager) joinRoomProcess(roomID int) {
	user, _ := clientsMap.Get(info.FromUser)
	joinedRoom := info.doRoomJoin(roomID, &user)

	roomNotificationQueue <- BroadcastRequest{
		broadcast: OutgoingBroadcast{
			Type:      "room_update",
			Message:   fmt.Sprintf("%s joined", user.username),
			RoomName:  joinedRoom.Name,
			RoomID:    joinedRoom.ID,
			RoomUsers: getRoomMemberNames(joinedRoom.ID),
		},
		receivers: roomStorage.GetRoomMemberConnections(joinedRoom.ID),
	}
}

// Data logic for joining a room
func (info *CommunicationsManager) doRoomJoin(roomID int, user *User) Room {
	if user.currentRoom >= 0 {
		roomStorage.RemoveMember(user.currentRoom, info.FromUser)
	} else {
		lobby.Delete(info.FromUser)
	}
	joinedRoom := roomStorage.AddMember(roomID, info.FromUser)
	*user = clientsMap.AddUserToGroup(info.FromUser, roomID)
	return joinedRoom
}

// Handle incoming request to register a username
func (info *CommunicationsManager) register() {
	user, ok := clientsMap.Get(info.FromUser)
	if ok {
		initUserData(&user, info.IncomingMessage.Username)
		clientsMap.Set(info.FromUser, user)
		lobby.Set(info.FromUser)
		reply := OutgoingBroadcast{
			Type:     "room_list",
			RoomList: roomStorage.GetAllProxied(),
		}
		sendBroadcast(user.connectionKey, reply)
	}
}

// Handle incoming request to leave the room
func (info *CommunicationsManager) leaveRoom() {
	user, _ := clientsMap.Get(info.FromUser)
	currentRoomID := user.currentRoom
	leaveRoom(info.FromUser)
	lobby.Set(info.FromUser)

	reply := OutgoingBroadcast{
		Type:     "room_list",
		RoomList: roomStorage.GetAllProxied(),
	}
	user, _ = clientsMap.Get(info.FromUser)
	sendBroadcast(user.connectionKey, reply)

	roomNotificationQueue <- BroadcastRequest{
		broadcast: OutgoingBroadcast{
			Type:      "room_update",
			Message:   fmt.Sprintf("%s left", user.username),
			RoomUsers: getRoomMemberNames(currentRoomID),
		},
		receivers: roomStorage.GetRoomMemberConnections(currentRoomID),
	}

}
