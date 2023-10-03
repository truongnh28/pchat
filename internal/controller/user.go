package controller

import (
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/internal/service"
	chat_app "chat-app/proto/chat-app"
	"context"
)

type user struct {
	userService service.UserService
}

func (u *user) SearchFriend(
	ctx context.Context,
	request *chat_app.SearchUserRequest,
) (resp *chat_app.SearchUserResponse, err error) {
	var (
		errCode        = common.OK
		userId, logger = helper.GetUserAndLogger(ctx)
		users          = make([]domain.User, 0)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.SearchUserResponse)
	users, errCode = u.userService.GetByUserName(ctx, request.GetUsername(), userId)
	if errCode != common.OK {
		logger.Errorln("get by user name failed")
		return
	}
	for _, i := range users {
		resp.UserInfos = append(resp.UserInfos, &chat_app.UserInfo{
			Username:    i.Username,
			PhoneNumber: i.PhoneNumber,
			Email:       i.Email,
			Status:      i.Status,
			UserId:      i.UserId,
			DateOfBirth: i.DateOfBirth.String(),
			Gender:      i.Gender.String(),
			Url:         i.Url,
		})
	}
	return
}

func (u *user) Search(
	ctx context.Context,
	request *chat_app.SearchUserRequest,
) (resp *chat_app.SearchUserResponse, err error) {
	var (
		errCode   = common.OK
		_, logger = helper.GetUserAndLogger(ctx)
		users     = make([]domain.User, 0)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.SearchUserResponse)
	users, errCode = u.userService.GetByUserName(ctx, request.GetUsername(), "")
	if errCode != common.OK {
		logger.Errorln("get by user name failed")
		return
	}
	for _, i := range users {
		resp.UserInfos = append(resp.UserInfos, &chat_app.UserInfo{
			Username:    i.Username,
			PhoneNumber: i.PhoneNumber,
			Email:       i.Email,
			Status:      i.Status,
			UserId:      i.UserId,
			DateOfBirth: i.DateOfBirth.String(),
			Gender:      i.Gender.String(),
			Url:         i.Url,
		})
	}
	return
}

func (u *user) Get(
	ctx context.Context,
	request *chat_app.GetUserRequest,
) (resp *chat_app.GetUserResponse, err error) {
	var (
		errCode      = common.OK
		User, logger = helper.GetUserAndLogger(ctx)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.GetUserResponse)

	logger.WithField("user", User)
	resp, errCode = u.userService.Get(ctx, request)
	if errCode != common.OK {
		logger.WithError(err).Error("get user list failed")
		return
	}
	return
}

func (u *user) Update(
	ctx context.Context,
	request *chat_app.UpdateUserRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode      = common.OK
		User, logger = helper.GetUserAndLogger(ctx)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.BasicResponse)

	logger.WithField("user", User)
	errCode = u.userService.UpdatePassword(ctx, request)
	if errCode != common.OK {
		logger.WithError(err).Error("get user list failed")
		return
	}
	return
}

func NewUser(
	userService service.UserService,
) chat_app.UserServer {
	return &user{
		userService: userService,
	}
}
