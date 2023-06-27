package firebase

import (
	"chat-app/internal/domain"
	"context"
	fb "firebase.google.com/go"
	"fmt"
	"github.com/whatvn/denny"
	"path/filepath"
	"sync"

	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type Firebase interface {
	SendToToken(
		ctx context.Context,
		notificationMessage domain.NotificationMessage,
		token string,
	) error
	SendToTopic(
		ctx context.Context,
		notificationMessage domain.NotificationMessage,
		topic string,
	) error
	SendMultiClient(
		ctx context.Context,
		notificationMessage domain.NotificationMessage,
		tokens []string,
	) error
}

type firebase struct {
	FirebaseApp     *fb.App
	MessagingClient *messaging.Client
}

var f *firebase
var firebaseOne sync.Once

func GetFirebase(certPath string) Firebase {
	ctx := context.Background()
	firebaseOne.Do(func() {

		serviceAccountKeyFilePath, err := filepath.Abs(certPath)
		if err != nil {
			panic("Unable to load serviceAccountKeys.json file")
		}
		opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
		app, err := fb.NewApp(ctx, nil, opt)
		if err != nil {
			panic(fmt.Errorf("unable to connect to firebase: %v", err.Error()))
		}

		client, err := app.Messaging(ctx)
		if err != nil {
			panic(fmt.Errorf("unable to connect to firebase messaging: %v", err.Error()))
		}

		f = &firebase{
			FirebaseApp:     app,
			MessagingClient: client,
		}
	})

	return f
}

func (f2 *firebase) SendToToken(
	ctx context.Context,
	notificationMessage domain.NotificationMessage,
	token string,
) error {
	logger := denny.GetLogger(ctx)
	client, err := f2.FirebaseApp.Messaging(ctx)
	if err != nil {
		logger.WithError(err).Errorf("error getting Messaging client: %v\n", err)
		return err
	}
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    notificationMessage.Title,
			Body:     notificationMessage.Body,
			ImageURL: notificationMessage.ImageURL,
		},
		Token: token,
	}

	_, err = client.Send(ctx, message)
	if err != nil {
		logger.WithError(err).Errorf("send noti fail: %v\n", err)
		return err
	}
	return nil
}

func (f2 *firebase) SendToTopic(
	ctx context.Context,
	notificationMessage domain.NotificationMessage,
	topic string,
) error {
	logger := denny.GetLogger(ctx)
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Topic: topic,
	}

	_, err := f2.MessagingClient.Send(ctx, message)
	if err != nil {
		logger.WithError(err).Errorf("send noti fail: %v\n", err)
		return err
	}

	return nil
}

func (f2 *firebase) SendMultiClient(
	ctx context.Context,
	notificationMessage domain.NotificationMessage,
	tokens []string,
) error {
	logger := denny.GetLogger(ctx)
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    notificationMessage.Title,
			Body:     notificationMessage.Body,
			ImageURL: notificationMessage.ImageURL,
		},
		Tokens: tokens,
	}

	_, err := f2.MessagingClient.SendMulticast(ctx, message)
	if err != nil {
		logger.WithError(err).Errorf("send noti fail: %v\n", err)
		return err
	}

	return nil
}
