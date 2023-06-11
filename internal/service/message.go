package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/pkg/repositories"
	"context"
)

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock
type MessageService interface {
	GetChatHistory(ctx context.Context, senderId, recipientId string) ([]*domain.ChatMessage, common.SubReturnCode)
	StoreNewChatMessages(ctx context.Context, message *domain.ChatMessage) error
}

type messageServiceImpl struct {
	messageRepository repositories.MessageRepository
}

func (m *messageServiceImpl) StoreNewChatMessages(ctx context.Context, message *domain.ChatMessage) error {
	return nil
}

func (m *messageServiceImpl) GetChatHistory(ctx context.Context, senderId, recipientId string) ([]*domain.ChatMessage, common.SubReturnCode) {
	//TODO implement me
	panic("implement me")
}

func NewMessageService(
	messageRepository repositories.MessageRepository,
) MessageService {
	return &messageServiceImpl{
		messageRepository: messageRepository,
	}
}
