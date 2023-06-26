package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/pkg/repositories"
	"context"
	"github.com/whatvn/denny"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock
type MessageService interface {
	GetChatHistory(
		ctx context.Context,
		isGroup bool,
		senderId, recipientId string,
		startTime, endTime time.Time,
	) ([]domain.ChatMessage, common.SubReturnCode)
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
	isGroup bool,
	senderId, recipientId string,
	startTime, endTime time.Time,
) ([]domain.ChatMessage, common.SubReturnCode) {
	var (
		logger       = denny.GetLogger(ctx)
		chatMessages = make([]domain.ChatMessage, 0)
	)
	resp, err := m.messageRepository.GetChatHistoryBetweenTwoUsers(
		ctx,
		isGroup,
		senderId,
		recipientId,
		startTime,
		endTime,
	)
	if err != nil {
		logger.WithError(err).Errorln("get chat history fail: ", err)
		return chatMessages, common.SystemError
	}
	return resp, common.OK
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
