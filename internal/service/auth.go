package service

import (
	"chat-app/config"
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/internal/repositories"
	"chat-app/models"
	cache "chat-app/pkg/client/redis"
	"chat-app/pkg/utils/auth"
	chat_app "chat-app/proto/chat-app"
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/whatvn/denny"
	"time"
)

type authenServiceImpl struct {
	serverCache cache.IRedisClient
	jWTAuth     auth.JWTAuth
	accountRepo repositories.UserRepository
	authConfig  *config.AuthenticationConfig
}

type AuthService interface {
	Login(ctx context.Context, req *chat_app.LoginRequest) (*chat_app.LoginResponse, error)
	VerifyOpt(ctx context.Context, verifyOtpReq *chat_app.VerifyOtpRequest) (bool, error)
	Logout(ctx context.Context, req *chat_app.LogoutRequest) error
}

func (a *authenServiceImpl) Login(
	ctx context.Context,
	req *chat_app.LoginRequest,
) (*chat_app.LoginResponse, error) {
	var (
		account *domain.User
		logger  = denny.GetLogger(ctx)
		resp    = new(chat_app.LoginResponse)
	)
	acc, err := a.accountRepo.FindByPhoneNumber(ctx, req.GetPhoneNumber())
	if err != nil {
		logger.WithError(err).Errorln("FindByUserName err: ", err)
		return resp, common.InvalidAccount
	}

	if acc.PhoneNumber != "" && helper.CheckPasswordHash(req.GetPassword(), acc.Password) {
		if acc.Status == models.Blocked {
			logger.Errorln("account blocked", acc.UserName)
			return resp, common.BlockedAccount
		}
	} else {
		logger.Errorln("wrong login information")
		return resp, common.LoginInfoInvalid
	}

	key := fmt.Sprintf("%s:%s", common.PrefixLoginCode, req.GetPhoneNumber())
	err = a.serverCache.Get(ctx, key).Err()
	if err != nil && err != redis.Nil {
		logger.Errorln("Login GetCode err: ", err)
		return resp, common.LoginSystemError
	}

	account = &domain.User{
		UserId:      acc.UserId,
		Username:    acc.UserName,
		Email:       acc.Email,
		PhoneNumber: acc.PhoneNumber,
		Status:      string(acc.Status),
	}

	if err == redis.Nil {
		account.Code = uuid.New().String()
	}

	err = a.serverCache.Set(ctx, key, account.Code, time.Duration(a.authConfig.ExpiredTime)).Err()
	if err != nil {
		logger.Errorln("SetCode failed:", err)
		return resp, common.LoginSystemError
	}

	newAcc, _ := sonic.Marshal(account)
	logger.Infoln("Initialize access token")
	token, err := a.jWTAuth.InitializeToken(string(newAcc))
	if err != nil {
		logger.Errorln("InitializeToken failed:", err)
		return resp, common.LoginSystemError
	}
	resp.AccessToken = token
	resp.UserId = acc.UserId
	return resp, nil
}

func (a *authenServiceImpl) Logout(ctx context.Context, req *chat_app.LogoutRequest) error {
	//TODO implement me
	panic("implement me")
}

func (a *authenServiceImpl) VerifyOpt(
	ctx context.Context,
	req *chat_app.VerifyOtpRequest,
) (bool, error) {
	var (
		logger = denny.GetLogger(ctx)
		err    = error(nil)
	)

	res := a.serverCache.Get(ctx, req.GetEmail())
	if res.Err() != nil {
		logger.WithError(res.Err()).Errorln("get otp from redis fail: ", res.Err())
		return false, res.Err()
	}

	if res.Val() != req.GetOtp() {
		logger.Errorln("wrong otp")
		return false, errors.New("wrong otp")
	}
	_, err = a.accountRepo.UpdateStatus(
		ctx,
		req.GetEmail(),
		models.Active,
	)
	if err != nil {
		logger.WithError(err).Errorln("update status fail")
		return false, err
	}
	return true, err
}

func NewAuthenService(
	jWTAuth auth.JWTAuth,
	serverCache cache.IRedisClient,
	accountRepo repositories.UserRepository,
	authConfig *config.AuthenticationConfig,
) AuthService {
	return &authenServiceImpl{
		serverCache: serverCache,
		jWTAuth:     jWTAuth,
		accountRepo: accountRepo,
		authConfig:  authConfig,
	}
}
