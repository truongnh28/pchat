package controller

import (
	"chat-app/config"
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/service"
	chat_app "chat-app/proto/chat-app"
	"context"
)

type account struct {
	accountService service.AccountService
	mailService    service.MailService
}

func (a account) Get(
	ctx context.Context,
	request *chat_app.GetAccountRequest,
) (resp *chat_app.GetAccountResponse, err error) {
	var (
		errCode         = common.OK
		account, logger = helper.GetAccountAndLogger(ctx)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.GetAccountResponse)

	logger.WithField("account", account)
	resp, errCode = a.accountService.Get(ctx, request)
	if errCode != common.OK {
		logger.WithError(err).Error("get role list failed")
		return
	}
	return
}

func (a account) Update(
	ctx context.Context,
	request *chat_app.UpdateAccountRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode         = common.OK
		account, logger = helper.GetAccountAndLogger(ctx)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.BasicResponse)

	logger.WithField("account", account)
	errCode = a.accountService.UpdatePassword(ctx, request)
	if errCode != common.OK {
		logger.WithError(err).Error("get role list failed")
		return
	}
	return
}

func (a account) Create(
	ctx context.Context,
	request *chat_app.CreateAccountRequest,
) (resp *chat_app.CreateAccountResponse, err error) {
	var (
		errCode         = common.OK
		account, logger = helper.GetAccountAndLogger(ctx)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.CreateAccountResponse)

	logger.WithField("account", account)

	passwordHash, err := helper.HashPassword(request.Password)
	if err != nil {
		errCode = common.SystemError
		return
	}
	opt, err := helper.GenOtp(config.GetAppConfig().Authentication.SecretKey)
	if err != nil {
		errCode = common.SystemError
		return
	}
	request.Password = passwordHash
	resp, errCode = a.accountService.Create(ctx, request)
	if errCode != common.OK {
		logger.WithError(err).Error("get role list failed")
		return
	}

	err = a.mailService.SendOpt([]string{request.Email}, opt)
	if err != nil {
		errCode = common.SystemError
		return
	}

	return
}

func NewAccount(
	accountService service.AccountService,
	mailService service.MailService,
) chat_app.AccountServer {
	return &account{
		accountService: accountService,
		mailService:    mailService,
	}
}
