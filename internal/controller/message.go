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
}

func (m *message) CreateMessage(
	ctx context.Context,
	request *chat_app.EmptyRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode        = common.OK
		userId, logger = helper.GetAccountAndLogger(ctx)
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
	roomId := httpCtx.Request.FormValue("room_id")
	if roomId == "" {
		errCode = common.InvalidRequest
		logger.WithError(errors.New("room_id is valid"))
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

	if text != "" {
		m.socketService.EmitNewMessage(ctx, roomId, &domain.ChatMessage{
			SenderID:    userId,
			RecipientID: roomId,
			Message:     text,
			Time:        time.Now(),
		})
	}

	if err == nil {
		uploadRes, errCode = m.mediaService.Upload(domain.UploadIn{
			FileName: fileHeader.Filename,
			FileData: file,
		})
		m.socketService.EmitNewMessage(ctx, roomId, &domain.ChatMessage{
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

	chatHistory, returnCode := m.messageService.GetChatHistory(ctx, req.SenderId, req.RecipientId)
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

func NewMessage(
	messageService service.MessageService,
	mediaService service.MediaService,
	socketService service.SocketService,
) chat_app.MessageServer {
	return &message{
		messageService: messageService,
		mediaService:   mediaService,
		socketService:  socketService,
	}
}
