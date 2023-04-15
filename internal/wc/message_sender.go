package wc

// EmitToSpecificClient will emit the socket event to specific socket user(sendChannel message to single user)
// Based on the userID provided, it sends socket Payload to that user, whos userId matches with the userId of
// users stored in socket server.
func EmitToSpecificClient(hub *Hub, payload SocketEventStruct, userId string) {
	for client := range hub.ConnectedClients {
		if client.userId == userId {
			for {
				select {
				case client.sendChannel <- payload:
				default:
					close(client.sendChannel)
					delete(hub.ConnectedClients, client)
				}
			}
		}
	}
}

// BroadcastSocketEventToAllClient will emit the socket events to all socket users
// function will sendChannel the socket events to all the users connected to Socket Server by using a ←sendChannel channel.
func BroadcastSocketEventToAllClient(hub *Hub, payload SocketEventStruct) {
	for client := range hub.ConnectedClients {
		select {
		case client.sendChannel <- payload:
		default:
			close(client.sendChannel)
			delete(hub.ConnectedClients, client)
		}
	}
}

func BroadcastSocketEventToAllClientExceptMe(hub *Hub, payload SocketEventStruct, myUserId string) {
	for client := range hub.ConnectedClients {
		if client.userId != myUserId {
			select {
			case client.sendChannel <- payload:
			default:
				close(client.sendChannel)
				delete(hub.ConnectedClients, client)
			}
		}
	}
}
