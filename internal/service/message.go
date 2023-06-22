package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/pkg/repositories"
	"context"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock
type MessageService interface {
	GetChatHistory(
		ctx context.Context,
		senderId, recipientId string,
		startTime, endTime time.Time,
	) ([]*domain.ChatMessage, common.SubReturnCode)
	CreateMessages(
		ctx context.Context,
		roomId string,
		message *domain.ChatMessage,
	) common.SubReturnCode
}

type messageServiceImpl struct {
	messageRepository repositories.MessageRepository
	socketService     SocketService
}

func (m *messageServiceImpl) CreateMessages(
	ctx context.Context,
	roomId string,
	message *domain.ChatMessage,
) common.SubReturnCode {
	// send message to socket
	m.socketService.EmitNewMessage(ctx, roomId, message)

	// store message
	// TODO: using kafka to store message
	err := m.messageRepository.StoreNewChatMessages(ctx, message)
	if err != nil {
		return common.SystemError
	}
	return common.OK
}

func (m *messageServiceImpl) GetChatHistory(
	ctx context.Context,
	senderId, recipientId string,
	startTime, endTime time.Time,
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
