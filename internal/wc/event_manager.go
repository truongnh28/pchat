package wc

import (
	"chat-app/internal/domain"
)

// HandleUserRegisterEvent will handle the join event for new socket users
func HandleUserRegisterEvent(hub *Hub, client *Client) {
	hub.ConnectedClients[client] = true
	handleSocketPayloadEvent(client, &SocketEventStruct{
		EventName:    domain.JoinEvent,
		EventPayload: client.userId,
	})
}

// HandleUserDisconnectEvent will handle the disconnect event for new socket users
func HandleUserDisconnectEvent(hub *Hub, client *Client) {
	_, ok := hub.ConnectedClients[client]
	if ok {
		delete(hub.ConnectedClients, client)
		close(client.sendChannel)
		handleSocketPayloadEvent(client, &SocketEventStruct{
			EventName:    domain.DisconnectEvent,
			EventPayload: client.userId,
		})
	}
}

type chatlistResponseStruct struct {
	Type     string      `json:"type"`
	Chatlist interface{} `json:"chatlist"`
}

func handleSocketPayloadEvent(client *Client, eventStruct *SocketEventStruct) {
	switch eventStruct.EventName {
	case domain.JoinEvent:
		//TODO: handle log
		userID := (eventStruct.EventPayload).(string)
		userDetail, err := client.hub.UserRepo.GetByUserId(userID)
		if err != nil {
			// TODO: handle log
			break
		}
		newUserOnline := SocketEventStruct{
			EventName: "new-user-joined",
			EventPayload: domain.UserDetail{
				ID:       userID,
				Username: userDetail.Username,
				Online:   userDetail.Online,
			},
		}
		BroadcastSocketEventToAllClientExceptMe(client.hub, newUserOnline, userID)

		allOnlineUsers := SocketEventStruct{
			EventName: "chatlist-response",
			EventPayload: chatlistResponseStruct{
				Type:     "my-chat-list",
				Chatlist: client.hub.UserRepo.GetAllOnlineUsers(userID),
			},
		}
		EmitToSpecificClient(client.hub, allOnlineUsers, userID)
	case domain.DisconnectEvent:
		if eventStruct.EventPayload != nil {
			userId := (eventStruct.EventPayload).(string)
			userDetail, err := client.hub.UserRepo.GetByUserId(userId)
			if err != nil {
				// TODO: handle log
				break
			}
			_ = client.hub.UserRepo.UpdateUserOnlineStatusByUserID(userId, false)
			BroadcastSocketEventToAllClient(client.hub, SocketEventStruct{
				EventName: "chatlist-response",
				EventPayload: chatlistResponseStruct{
					Type: "user-disconnect",
					Chatlist: domain.UserDetail{
						ID:       userId,
						Username: userDetail.Username,
						Online:   false,
					},
				},
			})
		}
	case domain.MessageEvent:
		message := (eventStruct.EventPayload).(map[string]interface{})["message"].(string)
		senderId := (eventStruct.EventPayload).(map[string]interface{})["senderId"].(string)
		recipientId := (eventStruct.EventPayload).(map[string]interface{})["recipientId"].(string)
		if !(message != "" && senderId != "" && recipientId != "") {
			break
		}
		messagePacket := domain.ChatMessage{
			SenderID:    senderId,
			RecipientID: recipientId,
			Message:     message,
		}
		client.hub.MessageRepo.StoreNewChatMessages(&messagePacket)
		allOnlientUser := SocketEventStruct{
			EventName:    "message-response",
			EventPayload: messagePacket,
		}
		EmitToSpecificClient(client.hub, allOnlientUser, recipientId)
	}
}
