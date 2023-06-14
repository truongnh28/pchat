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
)

type message struct {
	messageService service.MessageService
	mediaService   service.MediaService
}

func (m *message) CreateMessage(
	ctx context.Context,
	request *chat_app.EmptyRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode   = common.OK
		_, logger = helper.GetAccountAndLogger(ctx)
		ok        = false
		httpCtx   *denny.Context
		uploadRes *uploader.UploadResult
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
	file, fileHeader, err := httpCtx.Request.FormFile("file")
	if err != nil {
		errCode = common.InvalidRequest
		logger.Errorln("Get file from request err: ", err)
		return
	}
	uploadRes, errCode = m.mediaService.Upload(domain.UploadIn{
		FileName: fileHeader.Filename,
		FileData: file,
	})
	// text := httpCtx.Request.FormValue("text")

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

	if req.SenderId == "" {
		err = errors.New("sender id isn't valid")
		subReturnCode = common.InvalidRequest
		return
	}
	if req.RecipientId == "" {
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
) chat_app.MessageServer {
	return &message{
		messageService: messageService,
		mediaService:   mediaService,
	}
}
