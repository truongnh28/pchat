package controller

import (
	"chat-app/internal/common"
	"chat-app/internal/service"
	chat_app "chat-app/proto/chat-app"
	"context"
	"github.com/whatvn/denny"
)

type user struct {
	userService service.UserService
}

func NewUser(
	userService service.UserService,
) chat_app.UserServer {
	return &user{
		userService: userService,
	}
}

func (u *user) GetUserDetail(ctx context.Context, req *chat_app.UserDetailRequest) (
	resp *chat_app.UserDetailResponse,
	err error,
) {
	var (
		subReturnCode = common.OK
		logger        = denny.GetLogger(ctx).WithField("user", req)
	)

	defer func() {
		if err != nil {
			logger.WithError(err).Error("get user detail request failed")
		}
		buildResponse(subReturnCode, resp)
		err = nil
	}()

	resp = new(chat_app.UserDetailResponse)
	//TODO: handle
	return
}
