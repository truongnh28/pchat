package controller

import (
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/internal/service"
	chat_app "chat-app/proto/chat-app"
	"context"
	"errors"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/whatvn/denny"
	"net/http"
	"time"
)

type message struct {
	messageService service.MessageService
	mediaService   service.MediaService
	socketService  service.SocketService
	userService    service.UserService
	roomService    service.RoomService
}

func (m *message) CreateMessage(
	ctx context.Context,
	request *chat_app.EmptyRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode        = common.OK
		userId, logger = helper.GetUserAndLogger(ctx)
		ok             = false
		httpCtx        *denny.Context
		uploadRes      *uploader.UploadResult
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.BasicResponse)
	httpCtx, ok = ctx.(*denny.Context)
	if !ok {
		errCode = common.SystemError
		logger.WithError(errors.New("get httpCtx fail"))
		return
	}
	roomId := httpCtx.Request.FormValue("recipient_id")
	if roomId == "" {
		errCode = common.InvalidRequest
		logger.WithError(errors.New("recipient_id is valid"))
		return
	}
	err = httpCtx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		errCode = common.InvalidRequest
		logger.WithError(errors.New("cloud not parse form data"))
		return
	}
	text := httpCtx.Request.FormValue("text")
	file, fileHeader, err := httpCtx.Request.FormFile("file")
	if err == http.ErrMissingFile && text == "" {
		errCode = common.InvalidRequest
		logger.Errorln("Get file from request err: ", err)
		return
	}
	isValid := m.validateRoom(ctx, roomId)
	if !isValid {
		errCode = common.InvalidRequest
		logger.Errorln("Get room request err: ", err)
		return
	}

	if text != "" {
		errCode = m.messageService.CreateMessages(ctx, roomId, &domain.ChatMessage{
			SenderID:    userId,
			RecipientID: roomId,
			Message:     text,
			Time:        time.Now(),
		})
	}

	if err == nil {
		uploadRes, errCode = m.mediaService.Push(domain.UploadIn{
			FileName: fileHeader.Filename,
			FileData: file,
		})
		errCode = m.messageService.CreateMessages(ctx, roomId, &domain.ChatMessage{
			SenderID:    userId,
			RecipientID: roomId,
			Message:     uploadRes.SecureURL,
			Time:        time.Now(),
		})
	}

	return
}

func (m *message) GetChatHistory(ctx context.Context, req *chat_app.ChatHistoryRequest) (
	resp *chat_app.ChatHistoryResponse,
	err error,
) {
	var (
		subReturnCode = common.OK
		logger        = denny.GetLogger(ctx).WithField("message", req)
		startTime     time.Time
		endTime       time.Time
	)

	defer func() {
		if err != nil {
			logger.WithError(err).Error("get chat history request failed")
		}
		buildResponse(subReturnCode, resp)
		err = nil
	}()

	resp = new(chat_app.ChatHistoryResponse)

	if req.GetSenderId() == "" {
		err = errors.New("sender id isn't valid")
		subReturnCode = common.InvalidRequest
		return
	}
	if req.GetRecipientId() == "" {
		err = errors.New("recipient id isn't valid")
		subReturnCode = common.InvalidRequest
		return
	}
	startTime, err = helper.ParseClientTime(req.GetStartTime(), helper.APIClientDateTimeFormat)
	if err != nil {
		subReturnCode = common.SystemError
		logger.WithError(err)
		return
	}
	endTime, err = helper.ParseClientTime(req.GetEndTime(), helper.APIClientDateTimeFormat)
	if err != nil {
		subReturnCode = common.SystemError
		logger.WithError(err)
		return
	}
	chatHistory, returnCode := m.messageService.GetChatHistory(
		ctx,
		req.SenderId,
		req.RecipientId,
		startTime,
		endTime,
	)
	if returnCode != common.OK {
		err = fmt.Errorf("get chat history failed")
		subReturnCode = returnCode
		return
	}

	respChatHistory := make([]*chat_app.ChatMessage, 0)
	for _, it := range chatHistory {
		respChatHistory = append(respChatHistory, &chat_app.ChatMessage{
			SenderId:    it.SenderID,
			RecipientId: it.RecipientID,
			Message:     it.Message,
			Time:        it.Time.String(),
		})
	}
	resp.ChatHistory = respChatHistory
	return
}

func (m *message) validateRoom(ctx context.Context, roomId string) bool {
	var (
		logger   = denny.GetLogger(ctx)
		errUser  = true
		errGroup = true
	)
	// if send user to user: room_id == user_id_recipient
	// if send user to group: room_id == group_id
	_, errCode := m.userService.GetByUserId(ctx, roomId)
	if errCode != common.OK {
		logger.Errorln("check user fail: ", errCode)
		errUser = false
	}

	_, errCode = m.roomService.Get(ctx, domain.Room{
		GroupId: roomId,
	})
	if errCode != common.OK {
		logger.Errorln("check group fail: ", errCode)
		errGroup = false
	}
	return errUser || errGroup
}

func NewMessage(
	messageService service.MessageService,
	mediaService service.MediaService,
	socketService service.SocketService,
	userService service.UserService,
	roomService service.RoomService,
) chat_app.MessageServer {
	return &message{
		messageService: messageService,
		mediaService:   mediaService,
		socketService:  socketService,
		userService:    userService,
		roomService:    roomService,
	}
}
