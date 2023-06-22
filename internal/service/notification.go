package service

import (
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/pkg/client/firebase"
)

type NotificationService interface {
	Push(in domain.Notification) common.SubReturnCode
}

func NewNotificationService(fb firebase.Firebase, userService UserService) NotificationService {
	return &notificationServiceImpl{
		fb:          fb,
		userService: userService,
	}
}

type notificationServiceImpl struct {
	fb          firebase.Firebase
	userService UserService
}

func (m notificationServiceImpl) Push(
	in domain.Notification,
) common.SubReturnCode {

	return common.OK
}
