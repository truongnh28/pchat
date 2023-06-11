package repositories

import (
	"chat-app/internal/domain"
	"context"
	"github.com/whatvn/denny"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MessageRepository interface {
	StoreNewChatMessages(ctx context.Context, chatMessage *domain.ChatMessage) error
	GetChatHistoryBetweenTwoUsers(
		ctx context.Context,
		senderId string,
		recipientId string,
		timeStart, timeEnd time.Time,
	) ([]domain.ChatMessage, error)
}

type messageRepositoryImpl struct {
	database *mgo.Session
}

func (m *messageRepositoryImpl) StoreNewChatMessages(
	ctx context.Context,
	chatMessage *domain.ChatMessage,
) error {
	err := m.database.DB("chat-message").C("messages").Insert(chatMessage)
	if err != nil {
		denny.GetLogger(ctx).Error("StoreNewChatMessages err: ", err)
	}
	return err
}

func (m *messageRepositoryImpl) GetChatHistoryBetweenTwoUsers(
	ctx context.Context,
	senderId string,
	recipientId string,
	timeStart, timeEnd time.Time,
) ([]domain.ChatMessage, error) {
	collection := m.database.DB("chat-message").C("messages")
	var timeRange = bson.M{
		"$gte": timeStart,
		"$lte": timeEnd,
	}

	var filter = bson.M{
		"senderID":    senderId,
		"recipientID": recipientId,
		"time":        timeRange,
	}

	var results []domain.ChatMessage

	err := collection.Find(filter).All(&results)
	if err != nil {
		denny.GetLogger(ctx).Error("GetChatHistoryBetweenTwoUsers err: ", err)
	}
	return results, nil
}

func NewMessageRepository(database *mgo.Session) MessageRepository {
	return &messageRepositoryImpl{
		database: database,
	}
}
