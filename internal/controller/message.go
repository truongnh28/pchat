package controller

import (
	"chat-app/internal/common"
	"chat-app/internal/service"
	chat_app "chat-app/proto/chat-app"
	"context"
	"errors"
	"fmt"
	"github.com/whatvn/denny"
)

type message struct {
	messageService service.MessageService
}

func NewMessage(
	messageService service.MessageService,
) chat_app.ChatServer {
	return &message{
		messageService: messageService,
	}
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
