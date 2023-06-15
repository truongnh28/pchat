package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/pkg/repositories"
	"context"
)

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock
type MessageService interface {
	GetChatHistory(
		ctx context.Context,
		senderId, recipientId string,
	) ([]*domain.ChatMessage, common.SubReturnCode)
	CreateMessages(ctx context.Context, message *domain.ChatMessage) error
}

type messageServiceImpl struct {
	messageRepository repositories.MessageRepository
	socketService     SocketService
}

func (m *messageServiceImpl) CreateMessages(
	ctx context.Context,
	message *domain.ChatMessage,
) error {
	// send message to socket
	m.socketService.EmitNewMessage(ctx, message.RecipientID, message)
	// store message
	// TODO: using kafka
	return nil
}

func (m *messageServiceImpl) GetChatHistory(
	ctx context.Context,
	senderId, recipientId string,
) ([]*domain.ChatMessage, common.SubReturnCode) {
	//TODO implement me
	panic("implement me")
}

func NewMessageService(
	messageRepository repositories.MessageRepository,
	socketService SocketService,
) MessageService {
	return &messageServiceImpl{
		messageRepository: messageRepository,
		socketService:     socketService,
	}
}
