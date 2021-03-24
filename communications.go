package main

import "fmt"

type CommunicationsManager struct {
	incomingMessage IncomingMessage
	fromUser UserUUID
}

func (info *CommunicationsManager) sendMsg () {
	user, ok := clientsMap.Get(info.fromUser)
	if ok {
		broadcast <- BroadcastRequest{
			broadcast: OutgoingBroadcast{
				Type:      "message",
				Message:  info.incomingMessage.Message,
				Username: user.username,
				RoomID:    user.currentRoom,
			},
			receivers: roomStorage.GetRoomMemberConnections(user.currentRoom),
		}
	}
}

func (info *CommunicationsManager) createRoom () {

	// Init room and register in the storage
	createdRoom := createRoom(info.incomingMessage.RoomName, info.incomingMessage.Lat, info.incomingMessage.Lng)

	// Add user to room
	info.joinRoomProcess(createdRoom.ID)

	// Notify lobby of new room
	roomNotificationQueue <- BroadcastRequest{
		broadcast: OutgoingBroadcast{
			Type:      "room_list",
			RoomList: 	roomStorage.GetAllProxied(),
		},
		receivers: getLobbyUsersConnections(),
	}
}

// Handle the incoming request to join a room
func (info *CommunicationsManager) joinRoomProcess (roomID int) {
	user, _ := clientsMap.Get(info.fromUser)
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
		roomStorage.RemoveMember(user.currentRoom, info.fromUser)
	} else {
		lobby.Delete(info.fromUser)
	}
	joinedRoom := roomStorage.AddMember(roomID, info.fromUser)
	*user = clientsMap.AddUserToGroup(info.fromUser, roomID)
	return joinedRoom
}

// Handle incoming request to register a username
func (info *CommunicationsManager) register () {
	user, ok := clientsMap.Get(info.fromUser)
	if ok {
		initUserData(&user, info.incomingMessage.Username)
		clientsMap.Set(info.fromUser, user)
		lobby.Set(info.fromUser)
		reply := OutgoingBroadcast{
			Type:     "room_list",
			RoomList: roomStorage.GetAllProxied(),
		}
		sendBroadcast(user.connectionKey, reply)
	}
}

// Handle incoming request to leave the room
func (info *CommunicationsManager) leaveRoom () {
	user, _ := clientsMap.Get(info.fromUser)
	currentRoomID := user.currentRoom
	leaveRoom(info.fromUser)
	lobby.Set(info.fromUser)


	reply := OutgoingBroadcast{
		Type:     "room_list",
		RoomList:     roomStorage.GetAllProxied(),
	}
	user, _ = clientsMap.Get(info.fromUser)
	sendBroadcast(user.connectionKey, reply)

	roomNotificationQueue <- BroadcastRequest{
		broadcast: OutgoingBroadcast{
			Type:    "room_update",
			Message: fmt.Sprintf("%s left", user.username),
			RoomUsers: getRoomMemberNames(currentRoomID),
		},
		receivers: roomStorage.GetRoomMemberConnections(currentRoomID),
	}

}