package service

import (
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/internal/repositories"
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
	CreateMessage(
		ctx context.Context,
		roomId string,
		message *domain.ChatMessage,
	) common.SubReturnCode
	GetLastChatHistory(
		ctx context.Context,
		userId string,
	) ([]domain.RoomChatShortDetail, common.SubReturnCode)
	ValidateRoom(ctx context.Context, roomId string) (bool, bool)
}

type messageServiceImpl struct {
	messageRepository   repositories.MessageRepository
	socketService       SocketService
	notificationService NotificationService
	userService         UserService
	groupService        GroupService
	roomService         RoomService
}

func (m *messageServiceImpl) GetLastChatHistory(
	ctx context.Context,
	userId string,
) ([]domain.RoomChatShortDetail, common.SubReturnCode) {
	var (
		logger    = denny.GetLogger(ctx)
		resp      = make([]domain.RoomChatShortDetail, 0)
		roomName  = ""
		roomImage = ""
	)
	rooms, err := m.messageRepository.GetAllRoomHasMessage(ctx, userId)
	if err != nil {
		logger.WithError(err).Errorln("get all room has message fail: ", err)
		return resp, common.SystemError
	}
	for _, room := range rooms {
		chatMessage, err := m.messageRepository.GetLastChatHistory(
			ctx,
			userId,
			room,
		)
		if err != nil {
			logger.WithError(err).Errorln("get last chat history fail: ", err)
			return resp, common.SystemError
		}

		_, isGroup := m.ValidateRoom(ctx, room)
		if isGroup {
			group, errCode := m.groupService.Get(ctx, room)
			if errCode != common.OK {
				logger.WithError(err).Errorln("get group detail fail: ", err)
				continue
			}
			roomName = group.Name
			roomImage = group.ImageUrl
		} else {
			user, errCode := m.userService.GetByUserId(ctx, room)
			if errCode != common.OK {
				logger.WithError(err).Errorln("get group detail fail: ", err)
				continue
			}
			roomName = user.Username
			roomImage = user.Url
		}
		resp = append(resp, domain.RoomChatShortDetail{
			RoomName:  roomName,
			RoomImage: roomImage,
			RoomId:    room,
			IsGroup:   isGroup,
			Message:   chatMessage,
		})
	}

	return resp, common.OK
}

func (m *messageServiceImpl) CreateMessage(
	ctx context.Context,
	roomId string,
	message *domain.ChatMessage,
) common.SubReturnCode {
	_, logger := helper.GetUserAndLogger(ctx)
	// send message to socket
	m.socketService.EmitNewMessage(ctx, roomId, message)
	//userInfo, errCode := m.userService.GetByUserId(ctx, message.RecipientID)
	//if errCode != common.OK {
	//	logger.Errorf("get user by id fail")
	//	return common.SystemError
	//}
	//errCode = m.notificationService.Push(ctx, domain.Notification{
	//	UserId: userId,
	//	Message: domain.NotificationMessage{
	//		Title:    userInfo.Username,
	//		Body:     message.Message,
	//		ImageURL: userInfo.Url,
	//	},
	//})
	//if errCode != common.OK {
	//	logger.WithError(errors.New("push notification fail"))
	//}
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

func (m *messageServiceImpl) ValidateRoom(ctx context.Context, roomId string) (bool, bool) {
	var (
		logger   = denny.GetLogger(ctx)
		errUser  = true
		errGroup = true
		isGroup  = true
	)
	// if send user to user: room_id == user_id_recipient
	// if send user to group: room_id == group_id
	_, errCode := m.userService.GetByUserId(ctx, roomId)
	if errCode != common.OK {
		logger.Errorln("check user fail: ", errCode)
		errUser = false
	}
	if errCode == common.OK {
		isGroup = false
	}
	_, errCode = m.roomService.Get(ctx, domain.Room{
		GroupId: roomId,
	})
	if errCode != common.OK {
		logger.Errorln("check group fail: ", errCode)
		errGroup = false
	}
	return errUser || errGroup, isGroup
}

func NewMessageService(
	messageRepository repositories.MessageRepository,
	socketService SocketService,
	notificationService NotificationService,
	roomService RoomService,
	userService UserService,
	groupService GroupService,
) MessageService {
	return &messageServiceImpl{
		messageRepository:   messageRepository,
		socketService:       socketService,
		notificationService: notificationService,
		roomService:         roomService,
		userService:         userService,
		groupService:        groupService,
	}
}
