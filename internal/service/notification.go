package service

import (
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/pkg/client/firebase"
	"chat-app/pkg/client/redis"
	"context"
)

type NotificationService interface {
	Push(cxt context.Context, in domain.Notification) common.SubReturnCode
}

func NewNotificationService(
	fb firebase.Firebase,
	userService UserService,
	redisCli redis.IRedisClient,
) NotificationService {
	return &notificationService{
		fb:          fb,
		userService: userService,
		redisCli:    redisCli,
	}
}

type notificationService struct {
	fb          firebase.Firebase
	userService UserService
	redisCli    redis.IRedisClient
}

func (m *notificationService) Push(
	cxt context.Context,
	in domain.Notification,
) common.SubReturnCode {
	var (
		userId, logger = helper.GetUserAndLogger(cxt)
		deviceTokens   = make([]string, 0)
	)
	redisResp := m.redisCli.SMembers(cxt, userId)
	if redisResp.Err() != nil {
		logger.WithError(redisResp.Err()).
			Errorln("get device token from redis fail: ", redisResp.Err())
		return common.SystemError
	}
	deviceTokens = redisResp.Val()
	if userId != in.UserId {
		currentUserDeviceToken := deviceTokens
		redisResp := m.redisCli.SMembers(cxt, in.UserId)
		if redisResp.Err() != nil {
			logger.WithError(redisResp.Err()).
				Errorln("get device token from redis fail: ", redisResp.Err())
			return common.SystemError
		}
		deviceTokens = helper.GetUniqueElements(redisResp.Val(), currentUserDeviceToken)
	}

	err := m.fb.SendMultiClient(cxt, in.Message, deviceTokens)
	if err != nil {
		logger.WithError(err).Errorln("push notification fail: ", err)
		return common.SystemError
	}
	return common.OK
}
