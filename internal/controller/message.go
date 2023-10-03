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
	"github.com/whatvn/denny"
	"net/http"
	"sort"
	"time"
)

type message struct {
	messageService service.MessageService
	mediaService   service.FileService
}

func (m *message) GetLastChatHistory(
	ctx context.Context,
	req *chat_app.EmptyRequest,
) (resp *chat_app.GetAllChatRoomResponse, err error) {
	var (
		errCode        = common.OK
		userId, logger = helper.GetUserAndLogger(ctx)
	)

	defer func() {
		if err != nil {
			logger.WithError(err).Error("get chat history request failed")
		}
		buildResponse(errCode, resp)
		err = nil
	}()

	resp = new(chat_app.GetAllChatRoomResponse)

	roomChatDetails, returnCode := m.messageService.GetLastChatHistory(
		ctx,
		userId,
	)
	if returnCode != common.OK {
		err = fmt.Errorf("get chat history failed")
		errCode = returnCode
		return
	}
	sort.Slice(roomChatDetails, func(i, j int) bool {
		return roomChatDetails[i].Message.Time.Before(roomChatDetails[j].Message.Time)
	})
	for _, room := range roomChatDetails {
		resp.Room = append(resp.Room, &chat_app.RoomShortDetail{
			RoomName: room.RoomName,
			RoomAvt:  room.RoomImage,
			RoomId:   room.RoomId,
			IsGroup:  room.IsGroup,
			ChatMessage: &chat_app.ChatMessage{
				SenderId:     room.Message.SenderID,
				RecipientId:  room.Message.RecipientID,
				Message:      room.Message.Message,
				Time:         room.Message.Time.String(),
				FileName:     room.Message.FileName,
				Height:       room.Message.Height,
				Width:        room.Message.Width,
				FileSize:     room.Message.FileSize,
				Url:          room.Message.URL,
				ResourceType: string(room.Message.Type),
			},
		})
	}
	return
}

func (m *message) CreateMessage(
	ctx context.Context,
	request *chat_app.EmptyRequest,
) (resp *chat_app.CreateMessageResponse, err error) {
	var (
		errCode        = common.OK
		userId, logger = helper.GetUserAndLogger(ctx)
		ok             = false
		httpCtx        *denny.Context
		uploadRes      *domain.File
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.CreateMessageResponse)
	httpCtx, ok = ctx.(*denny.Context)
	if !ok {
		errCode = common.SystemError
		logger.WithError(common.GetHttpCtxFail)
		return
	}
	roomId := httpCtx.Request.FormValue("recipient_id")
	if roomId == "" {
		errCode = common.InvalidRequest
		logger.WithError(common.FiledInvalid)
		return
	}
	err = httpCtx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		errCode = common.InvalidRequest
		logger.WithError(common.ParseDataFail)
		return
	}
	text := httpCtx.Request.FormValue("text")
	file, fileHeader, err := httpCtx.Request.FormFile("file")
	if err == http.ErrMissingFile && text == "" {
		errCode = common.InvalidRequest
		logger.Errorln("Get file from request err: ", err)
		return
	}
	isValid, _ := m.messageService.ValidateRoom(ctx, roomId)
	if !isValid {
		errCode = common.InvalidRequest
		logger.Errorln("Get room request err")
		return
	}

	if text != "" {
		errCode = m.messageService.CreateMessage(ctx, roomId, &domain.ChatMessage{
			SenderID:    userId,
			RecipientID: roomId,
			Message:     text,
			Time:        time.Now(),
			Type:        domain.MessageText,
		})
	}

	if err == nil {
		uploadRes, errCode = m.mediaService.Create(ctx, domain.UploadIn{
			FileName: fileHeader.Filename,
			FileData: file,
		})
		errCode = m.messageService.CreateMessage(ctx, roomId, &domain.ChatMessage{
			SenderID:    userId,
			RecipientID: roomId,
			Message:     "",
			Time:        time.Now(),
			FileName:    uploadRes.GetOriginalFileName(),
			Height:      uploadRes.GetHeight(),
			Width:       uploadRes.GetWidth(),
			FileSize:    uploadRes.GetFileSize(),
			URL:         uploadRes.GetSecureURL(),
			Type: domain.StringToMessageType(
				uploadRes.ResourceType,
				uploadRes.GetSecureURL(),
			),
		})
		resp.Message = &chat_app.ChatMessage{
			SenderId:    userId,
			RecipientId: roomId,
			Message:     "",
			Time:        time.Now().String(),
			FileName:    uploadRes.OriginalFilename,
			Height:      uploadRes.Height,
			Width:       uploadRes.Width,
			FileSize:    uploadRes.FileSize,
			Url:         uploadRes.GetSecureURL(),
			ResourceType: string(
				domain.StringToMessageType(uploadRes.ResourceType, uploadRes.GetSecureURL()),
			),
		}
	}

	return
}

func (m *message) GetChatHistory(ctx context.Context, req *chat_app.ChatHistoryRequest) (
	resp *chat_app.ChatHistoryResponse,
	err error,
) {
	var (
		errCode        = common.OK
		userId, logger = helper.GetUserAndLogger(ctx)
		startTime      time.Time
		endTime        time.Time
	)

	defer func() {
		if err != nil {
			logger.WithError(err).Error("get chat history request failed")
		}
		buildResponse(errCode, resp)
		err = nil
	}()

	resp = new(chat_app.ChatHistoryResponse)

	if req.GetRecipientId() == "" {
		err = errors.New("recipient id isn't valid")
		errCode = common.InvalidRequest
		return
	}
	startTime, err = helper.ParseClientTime(req.GetStartTime(), helper.APIClientDateTimeFormat)
	if err != nil {
		errCode = common.SystemError
		logger.WithError(err)
		return
	}
	endTime, err = helper.ParseClientTime(req.GetEndTime(), helper.APIClientDateTimeFormat)
	if err != nil {
		errCode = common.SystemError
		logger.WithError(err)
		return
	}
	isValid, isGroup := m.messageService.ValidateRoom(ctx, req.GetRecipientId())
	if !isValid {
		errCode = common.InvalidRequest
		logger.Errorln("Get room request err: ", err)
		return
	}
	if userId == req.GetRecipientId() {
		errCode = common.InvalidRequest
		logger.Errorln("Get room request err: ", err)
		return
	}
	chatHistory, returnCode := m.messageService.GetChatHistory(
		ctx,
		isGroup,
		userId,
		req.RecipientId,
		startTime,
		endTime,
	)
	if returnCode != common.OK {
		err = fmt.Errorf("get chat history failed")
		errCode = returnCode
		return
	}

	respChatHistory := make([]*chat_app.ChatMessage, 0)
	for _, it := range chatHistory {
		respChatHistory = append(respChatHistory, &chat_app.ChatMessage{
			SenderId:     it.SenderID,
			RecipientId:  it.RecipientID,
			Message:      it.Message,
			Time:         it.Time.String(),
			FileName:     it.FileName,
			Height:       it.Height,
			Width:        it.Width,
			FileSize:     it.FileSize,
			Url:          it.URL,
			ResourceType: string(it.Type),
		})
	}
	resp.ChatHistory = respChatHistory
	return
}

func NewMessage(
	messageService service.MessageService,
	mediaService service.FileService,
) chat_app.MessageServer {
	return &message{
		messageService: messageService,
		mediaService:   mediaService,
	}
}
