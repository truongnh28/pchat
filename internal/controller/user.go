package controller

import (
	"chat-app/config"
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/service"
	cache "chat-app/pkg/client/redis"
	chat_app "chat-app/proto/chat-app"
	"context"
	"time"
)

type user struct {
	userService service.UserService
	mailService service.MailService
	redisCli    cache.IRedisClient
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
		logger.WithError(err).Error("get role list failed")
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
		logger.WithError(err).Error("get role list failed")
		return
	}
	return
}

func (u *user) Create(
	ctx context.Context,
	request *chat_app.CreateUserRequest,
) (resp *chat_app.CreateUserResponse, err error) {
	var (
		errCode   = common.OK
		_, logger = helper.GetUserAndLogger(ctx)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.CreateUserResponse)

	passwordHash, err := helper.HashPassword(request.Password)
	if err != nil {
		logger.WithError(err).Errorln("hash password fail")
		errCode = common.SystemError
		return
	}
	opt, err := helper.GenOtp(config.GetAppConfig().Authentication.SecretKey)
	if err != nil {
		logger.WithError(err).Errorln("gen otp fail")
		errCode = common.SystemError
		return
	}

	err = u.redisCli.Set(ctx, request.GetEmail(), opt, time.Minute*5).Err()
	if err != nil {
		logger.WithError(err).Errorln("gen otp fail")
		errCode = common.SystemError
		return
	}

	request.Password = passwordHash
	resp, errCode = u.userService.Create(ctx, request)
	if errCode != common.OK {
		logger.WithError(err).Error("create user fail")
		return
	}

	err = u.mailService.SendOpt([]string{request.GetEmail()}, opt)
	if err != nil {
		logger.WithError(err).Error("send otp fail")
		errCode = common.SystemError
		return
	}

	return
}

func NewUser(
	UserService service.UserService,
	mailService service.MailService,
	redisCli cache.IRedisClient,
) chat_app.UserServer {
	return &user{
		userService: UserService,
		mailService: mailService,
		redisCli:    redisCli,
	}
}
