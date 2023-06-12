package repositories

import (
	"chat-app/internal/domain"
	"context"
	"github.com/bytedance/sonic"
	"github.com/whatvn/denny"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
	senderId string,
	recipientId string,
	timeStart, timeEnd time.Time,
) ([]domain.ChatMessage, error) {
	collection := m.database.Collection("messages")
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
	findOptions := options.Find().SetSort(bson.D{{"time", 1}})
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		denny.GetLogger(ctx).Error("GetChatHistoryBetweenTwoUsers err: ", err)
	}
	// Iterate through the cursor and print the documents
	for cursor.Next(ctx) {
		var message domain.ChatMessage
		err := cursor.Decode(&message)
		jsonMessage, err := sonic.Marshal(message)
		if err != nil {
			log.Fatal(err)
		}
		err = sonic.Unmarshal(jsonMessage, &message)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
	}
	return results, nil
}

func NewMessageRepository(database *mongo.Database) MessageRepository {
	return &messageRepositoryImpl{
		database: database,
	}
}
