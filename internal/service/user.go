package service

import (
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/models"
	"chat-app/pkg/repositories"
	chat_app "chat-app/proto/chat-app"
	"context"
	"github.com/google/uuid"
	"github.com/whatvn/denny"
	"strings"
)

type UserService interface {
	Create(
		ctx context.Context,
		req *chat_app.RegisterRequest,
	) (*chat_app.RegisterResponse, common.SubReturnCode)
	UpdatePassword(
		ctx context.Context,
		req *chat_app.UpdateUserRequest,
	) common.SubReturnCode
	Get(
		ctx context.Context,
		req *chat_app.GetUserRequest,
	) (*chat_app.GetUserResponse, common.SubReturnCode)
	GetByUserId(
		ctx context.Context,
		userId string,
	) (domain.User, common.SubReturnCode)
}

type userService struct {
	userRepository repositories.UserRepository
	fileService    FileService
}

func (a *userService) GetByUserId(
	ctx context.Context,
	userId string,
) (domain.User, common.SubReturnCode) {
	var (
		err    = error(nil)
		logger = denny.GetLogger(ctx)
		acc    = &models.User{}
	)
	if userId == "" {
		return domain.User{}, common.SystemError
	}
	acc, err = a.userRepository.Get(ctx, userId)
	if err != nil {
		logger.WithError(err)
		return domain.User{}, common.SystemError
	}
	file, errCode := a.fileService.Get(ctx, acc.FileId)
	if errCode != common.OK {
		logger.Errorf("Find user file fail: %d", errCode)
		return domain.User{}, common.SystemError
	}
	return domain.User{
		UserId:      acc.UserId,
		Username:    acc.UserName,
		Email:       acc.Email,
		PhoneNumber: acc.PhoneNumber,
		Status:      string(acc.Status),
		Code:        "",
		DateOfBirth: acc.DateOfBirth,
		Gender:      acc.Gender,
		Url:         file.GetSecureURL(),
	}, common.OK
}

func (a *userService) Get(
	ctx context.Context,
	req *chat_app.GetUserRequest,
) (*chat_app.GetUserResponse, common.SubReturnCode) {
	var (
		acc    = &models.User{}
		err    = error(nil)
		resp   = &chat_app.GetUserResponse{}
		logger = denny.GetLogger(ctx)
	)
	if req.GetUsername() == "" && req.GetUserId() == "" {
		return resp, common.SystemError
	}
	if req.GetUserId() != "" {
		acc, err = a.userRepository.Get(ctx, req.GetUserId())
	} else {
		acc, err = a.userRepository.FindByUserName(ctx, req.GetUsername())
	}
	if err != nil {
		logger.WithError(err).Errorf("Find user fail: %s", err)
		return resp, common.SystemError
	}

	file, errCode := a.fileService.Get(ctx, acc.FileId)
	if errCode != common.OK {
		logger.Errorf("Find user file fail: %d", errCode)
		return resp, common.SystemError
	}
	resp.Info = &chat_app.UserInfo{
		Username:    acc.UserName,
		PhoneNumber: acc.PhoneNumber,
		Email:       acc.Email,
		Status:      string(acc.Status),
		UserId:      req.GetUserId(),
		DateOfBirth: helper.FromTimeToString(acc.DateOfBirth, helper.ApiClientDateFormat),
		Gender:      acc.Gender.String(),
		Url:         file.GetSecureURL(),
	}
	return resp, common.OK
}

func (a *userService) UpdatePassword(
	ctx context.Context,
	req *chat_app.UpdateUserRequest,
) common.SubReturnCode {
	var (
		logger = denny.GetLogger(ctx)
	)
	rowsAffected, err := a.userRepository.UpdatePassword(ctx, req)
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

func (a *userService) Create(
	ctx context.Context,
	req *chat_app.RegisterRequest,
) (*chat_app.RegisterResponse, common.SubReturnCode) {
	var (
		resp             = &chat_app.RegisterResponse{}
		logger           = denny.GetLogger(ctx)
		userId, _        = uuid.NewUUID()
		dateOfBirth, err = helper.ParseLocalTime(helper.ApiClientDateFormat, req.DateOfBirth)
	)
	if err != nil {
		logger.WithError(err).Errorln("convert time err: ", err)
		return resp, common.SystemError
	}
	acc := models.User{
		UserId:      userId.String(),
		UserName:    req.GetUsername(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Password:    req.GetPassword(),
		Status:      models.Blocked,
		DateOfBirth: dateOfBirth,
		Gender:      helper.ConvertToGenderType(req.Gender),
		FileId:      1,
	}
	err = a.userRepository.Validate(ctx, acc)
	if err != nil {
		errStr := "Create User service err: "
		if strings.Contains(err.Error(), "is exist") {
			errStr = "Invalid request: "
		}
		logger.Errorln(errStr, err)
		return resp, common.SystemError
	}
	err = a.userRepository.Create(ctx, acc)
	if err != nil {
		logger.Errorln("Create User service err: ", err)
		return resp, common.SystemError
	}
	resp.Info = &chat_app.UserInfo{
		UserId:      userId.String(),
		Username:    req.GetUsername(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Status:      string(models.Blocked),
	}
	return resp, common.OK
}

func NewUserService(
	userRepository repositories.UserRepository,
	fileService FileService,
) UserService {
	return &userService{
		userRepository: userRepository,
		fileService:    fileService,
	}
}
