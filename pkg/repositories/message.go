package repositories

import (
	"chat-app/internal/domain"
	"context"
	"github.com/bytedance/sonic"
	"github.com/whatvn/denny"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MessageRepository interface {
	StoreNewChatMessages(ctx context.Context, chatMessage *domain.ChatMessage) error
	GetChatHistoryBetweenTwoUsers(
		ctx context.Context,
		isGroup bool,
		senderId string,
		recipientId string,
		timeStart, timeEnd time.Time,
	) ([]domain.ChatMessage, error)
}

type messageRepositoryImpl struct {
	database *mongo.Database
}

func (m *messageRepositoryImpl) StoreNewChatMessages(
	ctx context.Context,
	chatMessage *domain.ChatMessage,
) error {
	_, err := m.database.Collection("messages").InsertOne(ctx, chatMessage)
	if err != nil {
		denny.GetLogger(ctx).Error("StoreNewChatMessages err: ", err)
	}
	return err
}

func (m *messageRepositoryImpl) GetChatHistoryBetweenTwoUsers(
	ctx context.Context,
	isGroup bool,
	senderId string,
	recipientId string,
	timeStart, timeEnd time.Time,
) ([]domain.ChatMessage, error) {
	var (
		logger     = denny.GetLogger(ctx)
		collection = m.database.Collection("messages")
	)
	var timeRange = bson.M{
		"$gte": timeStart,
		"$lte": timeEnd,
	}

	var filter = bson.M{
		"senderid":    senderId,
		"recipientid": recipientId,
		"time":        timeRange,
	}

	if isGroup {
		filter = bson.M{
			"recipientid": recipientId,
			"time":        timeRange,
		}
	}

	var results []domain.ChatMessage
	findOptions := options.Find().SetSort(bson.D{{"time", 1}})
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		logger.Error("GetChatHistoryBetweenTwoUsers err: ", err)
		return nil, err
	}
	for cursor.Next(ctx) {
		var message domain.ChatMessage
		err := cursor.Decode(&message)
		jsonMessage, err := sonic.Marshal(message)
		if err != nil {
			logger.Error("marshal message err: ", err)
			return nil, err
		}
		err = sonic.Unmarshal(jsonMessage, &message)
		if err != nil {
			logger.Error("unmarshal message err: ", err)
			return nil, err
		}
		results = append(results, message)
	}
	return results, nil
}

func NewMessageRepository(database *mongo.Database) MessageRepository {
	return &messageRepositoryImpl{
		database: database,
	}
}
