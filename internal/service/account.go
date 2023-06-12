package service

import (
	"chat-app/internal/common"
	"chat-app/models"
	"chat-app/pkg/repositories"
	chat_app "chat-app/proto/chat-app"
	"context"
	"github.com/golang/glog"
	"github.com/whatvn/denny"
	"strings"
)

type AccountService interface {
	Create(
		ctx context.Context,
		req *chat_app.CreateAccountRequest,
	) (*chat_app.CreateAccountResponse, common.SubReturnCode)
	UpdatePassword(
		ctx context.Context,
		req *chat_app.UpdateAccountRequest,
	) common.SubReturnCode
	Get(
		ctx context.Context,
		req *chat_app.GetAccountRequest,
	) (*chat_app.GetAccountResponse, common.SubReturnCode)
}

type accountService struct {
	accountRepository repositories.AccountRepository
}

func (a *accountService) Get(
	ctx context.Context,
	req *chat_app.GetAccountRequest,
) (*chat_app.GetAccountResponse, common.SubReturnCode) {
	var (
		acc    = &models.Account{}
		err    = error(nil)
		resp   = &chat_app.GetAccountResponse{}
		logger = denny.GetLogger(ctx)
	)
	if req.GetUsername() == "" && req.GetUserId() == 0 {
		return resp, common.SystemError
	}
	if req.GetUserId() != 0 {
		acc, err = a.accountRepository.Get(ctx, req.GetUserId())
	} else {
		acc, err = a.accountRepository.FindByUserName(ctx, req.GetUsername())
	}
	if err != nil {
		glog.Errorf("Find user fail: %s", err)
		logger.WithError(err)
		return resp, common.SystemError
	}
	resp.Info = &chat_app.AccountInfo{
		Username:    acc.UserName,
		PhoneNumber: acc.Email,
		Email:       acc.PhoneNumber,
		Status:      string(acc.Status),
	}
	return resp, common.OK
}

func (a *accountService) UpdatePassword(
	ctx context.Context,
	req *chat_app.UpdateAccountRequest,
) common.SubReturnCode {
	var (
		logger = denny.GetLogger(ctx)
	)
	rowsAffected, err := a.accountRepository.UpdatePassword(ctx, req)
	if err != nil {
		logger.WithError(err).Errorln("Update failed: ", err)
		return common.SystemError
	}

	if rowsAffected == 0 {
		logger.Errorln("user not exist", req.GetUsername())
		return common.SystemError
	}

	return common.OK
}

func (a *accountService) Create(
	ctx context.Context,
	req *chat_app.CreateAccountRequest,
) (*chat_app.CreateAccountResponse, common.SubReturnCode) {
	var (
		resp   = &chat_app.CreateAccountResponse{}
		logger = denny.GetLogger(ctx)
	)
	acc := models.Account{
		UserName:    req.GetUsername(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Password:    req.GetPassword(),
		Status:      models.Blocked,
	}
	err := a.accountRepository.Validate(ctx, acc)
	if err != nil {
		errStr := "Create Account service err: "
		if strings.Contains(err.Error(), "is exist") {
			errStr = "Invalid request: "
		}
		logger.Errorln(errStr, err)
		return resp, common.SystemError
	}
	err = a.accountRepository.Create(ctx, acc)
	if err != nil {
		logger.Errorln("Create Account service err: ", err)
		return resp, common.SystemError
	}
	resp.Info = &chat_app.AccountInfo{
		Username:    req.GetUsername(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Status:      string(models.Blocked),
	}
	return resp, common.OK
}

func NewAccountService(
	accountRepository repositories.AccountRepository,
) AccountService {
	return &accountService{
		accountRepository: accountRepository,
	}
}
