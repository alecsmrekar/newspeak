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
	info.joinRoom(createdRoom.ID)

	// Notify lobby of new room
	roomNotificationQueue <- BroadcastRequest{
		broadcast: OutgoingBroadcast{
			Type:      "room_list",
			RoomList: 	roomStorage.GetAllProxied(),
		},
		receivers: getLobbyUsersConnections(),
	}
}

func (info *CommunicationsManager) joinRoom (roomID int) {
	user, _ := clientsMap.Get(info.fromUser)

	if user.currentRoom >= 0 {
		roomStorage.RemoveMember(user.currentRoom, info.fromUser)
	} else {
		lobby.Delete(info.fromUser)
	}
	joinedRoom := roomStorage.AddMember(roomID, info.fromUser)
	user = clientsMap.AddUserToGroup(info.fromUser, roomID)
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

func (info *CommunicationsManager) register () {
	user, ok := clientsMap.Get(info.fromUser)
	if ok {
		strategy := &register{}
		strategy.update(&user, UserPayload{message: info.incomingMessage})
		clientsMap.Set(info.fromUser, user)
		lobby.Set(info.fromUser)
		reply := OutgoingBroadcast{
			Type:     "room_list",
			RoomList: roomStorage.GetAllProxied(),
		}
		sendBroadcast(user.connectionKey, reply)
	}
}

func (info *CommunicationsManager) leaveRoom () {
	user, _ := clientsMap.Get(info.fromUser)
	currentRoomID := user.currentRoom
	leaveRoom(info.fromUser)
	currentRoom, ok := roomStorage.GetRoom(currentRoomID)
	user, _ = clientsMap.Get(info.fromUser)
	lobby.Set(info.fromUser)
	reply := OutgoingBroadcast{
		Type:     "room_list",
		RoomList:     roomStorage.GetAllProxied(),
	}
	sendBroadcast(user.connectionKey, reply)
	if ok {
		roomNotificationQueue <- BroadcastRequest{
			broadcast: OutgoingBroadcast{
				Type:    "room_update",
				Message: fmt.Sprintf("%s left", user.username),
				RoomID:  currentRoomID,
				RoomName: currentRoom.Name,
				RoomUsers: getRoomMemberNames(currentRoomID),
			},
			receivers: roomStorage.GetRoomMemberConnections(currentRoomID),
		}
	}
}