package service

import (
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/pkg/repositories"
	"context"
	"errors"
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
	CreateMessage(
		ctx context.Context,
		roomId string,
		message *domain.ChatMessage,
	) common.SubReturnCode
}

type messageServiceImpl struct {
	messageRepository   repositories.MessageRepository
	socketService       SocketService
	notificationService NotificationService
	userService         UserService
	groupService        GroupService
}

func (m *messageServiceImpl) CreateMessage(
	ctx context.Context,
	roomId string,
	message *domain.ChatMessage,
) common.SubReturnCode {
	userId, logger := helper.GetUserAndLogger(ctx)
	// send message to socket
	m.socketService.EmitNewMessage(ctx, roomId, message)
	userInfo, errCode := m.userService.GetByUserId(ctx, message.RecipientID)
	if errCode != common.OK {
		logger.Errorf("get user by id fail")
		return common.SystemError
	}
	errCode = m.notificationService.Push(ctx, domain.Notification{
		UserId: userId,
		Message: domain.NotificationMessage{
			Title:    userInfo.Username,
			Body:     message.Message,
			ImageURL: userInfo.Url,
		},
	})
	if errCode != common.OK {
		logger.WithError(errors.New("push notification fail"))
	}
	// store message
	// TODO: using kafka to store message
	err := m.messageRepository.StoreNewChatMessages(ctx, message)
	if err != nil {
		logger.WithError(err).Errorln("store message fail: ", err)
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
	notificationService NotificationService,
) MessageService {
	return &messageServiceImpl{
		messageRepository:   messageRepository,
		socketService:       socketService,
		notificationService: notificationService,
	}
}
