package service

import (
	"chat-app/internal/domain"
	"chat-app/internal/ws"
	"context"
	"github.com/bytedance/sonic"
	"github.com/whatvn/denny"
)

// SocketService defines methods related emitting websockets events the service layer expects
// any repository it interacts with to implement
type SocketService interface {
	EmitNewMessage(ctx context.Context, roomId string, message *domain.ChatMessage)
	EmitEditMessage(room string, message *domain.ChatMessage)
	EmitDeleteMessage(room, messageId string)
	EmitNewRoom(ctx context.Context, roomId string, room *domain.Group)
	EmitEditRoom(ctx context.Context, roomId string, room *domain.Group)
	EmitDeleteRoom(ctx context.Context, roomId string)
	EmitAddMember(room string, member *domain.UserDetail)
	EmitRemoveMember(room, memberId string)
	EmitNewDMNotification(roomId string, userId string)
	EmitSendRequest(roomId string)
	EmitAddFriendRequest(roomId string, request *domain.FriendRequest)
	EmitAddFriend(user, member *domain.UserDetail)
	EmitRemoveFriend(userId, memberId string)
}

type socketService struct {
	hub *ws.Hub
}

func (s *socketService) EmitEditMessage(room string, message *domain.ChatMessage) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitDeleteMessage(room, messageId string) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitNewRoom(ctx context.Context, roomId string, room *domain.Group) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitEditRoom(ctx context.Context, roomId string, room *domain.Group) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitDeleteRoom(ctx context.Context, roomId string) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitAddMember(room string, member *domain.UserDetail) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitRemoveMember(room, memberId string) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitNewDMNotification(roomId string, userId string) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitSendRequest(roomId string) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitAddFriendRequest(roomId string, request *domain.FriendRequest) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitAddFriend(user, member *domain.UserDetail) {
	//TODO implement me
	panic("implement me")
}

func (s *socketService) EmitRemoveFriend(userId, memberId string) {
	//TODO implement me
	panic("implement me")
}

// NewSocketService is a factory function for
// initializing a SocketService with its repository layer dependencies
func NewSocketService(hub *ws.Hub) SocketService {
	return &socketService{
		hub: hub,
	}
}

func (s *socketService) EmitNewMessage(
	ctx context.Context,
	roomId string,
	message *domain.ChatMessage,
) {
	logger := denny.GetLogger(ctx)
	data, err := sonic.Marshal(ws.SocketMessage{
		Event:   ws.NewMessage,
		Payload: message,
	})

	if err != nil {
		logger.Printf("error marshalling response: %v\n", err)
	}

	s.hub.BroadcastToRoom(data, roomId)
}
