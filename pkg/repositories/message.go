package repositories

import (
	"chat-app/internal/domain"
	"gopkg.in/mgo.v2"
)

type MessageRepository interface {
	StoreNewChatMessages(chatMessage *domain.ChatMessage) bool
	GetChatHistoryBetweenTwoUsers(senderId string, recipientId string) []*domain.MessageConversation
}

type messageRepositoryImpl struct {
	database *mgo.Session
}

func (m *messageRepositoryImpl) StoreNewChatMessages(chatMessage *domain.ChatMessage) bool {
	//TODO implement me
	panic("implement me")
}

func (m *messageRepositoryImpl) GetChatHistoryBetweenTwoUsers(senderId string, recipientId string) []*domain.MessageConversation {
	//TODO implement me
	panic("implement me")
}

func NewMessageRepository(database *mgo.Session) MessageRepository {
	return &messageRepositoryImpl{
		database: database,
	}
}
