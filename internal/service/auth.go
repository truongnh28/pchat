package service

import (
	"chat-app/config"
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/domain"
	"chat-app/models"
	cache "chat-app/pkg/client/redis"
	"chat-app/pkg/repositories"
	"chat-app/pkg/utils/auth"
	chat_app "chat-app/proto/chat-app"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/whatvn/denny"
	"strconv"
	"time"
)

type authenServiceImpl struct {
	serverCache *cache.RedisClient
	jWTAuth     auth.JWTAuth
	accountRepo repositories.AccountRepository
	authConfig  *config.AuthenticationConfig
}

type AuthenService interface {
	Login(ctx context.Context, req *chat_app.LoginRequest) (*chat_app.LoginResponse, error)
	VerifyOpt(ctx context.Context, verifyOtpReq *chat_app.VerifyOtpRequest) (bool, error)
	Logout(ctx context.Context, req *chat_app.LogoutRequest) error
}

func (a *authenServiceImpl) Login(
	ctx context.Context,
	req *chat_app.LoginRequest,
) (*chat_app.LoginResponse, error) {
	var (
		account *domain.Account
		logger  = denny.GetLogger(ctx)
		resp    = new(chat_app.LoginResponse)
	)
	acc, err := a.accountRepo.FindByPhoneNumber(ctx, req.GetPhoneNumber())
	if err != nil {
		logger.WithError(err).Errorln("FindByUserName err: ", err)
		return nil, errors.New("account not valid")
	}

	if acc.PhoneNumber != "" && helper.CheckPasswordHash(req.GetPassword(), acc.Password) {
		if acc.Status == models.Blocked {
			logger.Errorln("account blocked", acc.UserName)
			return nil, errors.New("account has been blocked")
		}
	} else {
		logger.Errorln("wrong login information")
		return nil, errors.New("wrong login information")
	}

	key := fmt.Sprintf("%s:%s", common.PrefixLoginCode, req.GetPhoneNumber())
	err = a.serverCache.Get(ctx, key).Err()
	if err != nil && err != redis.Nil {
		logger.Errorln("Login GetCode err: ", err)
		return nil, errors.New("system error")
	}

	account = &domain.Account{
		UserId:      uint64(acc.ID),
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
		return nil, errors.New("system error")
	}

	newAcc, _ := json.Marshal(account)
	logger.Infoln("Initialize access token")
	token, err := a.jWTAuth.InitializeToken(string(newAcc))
	if err != nil {
		logger.Errorln("InitializeToken failed:", err)
		return nil, errors.New("system error")
	}
	resp.AccessToken = token
	resp.UserId = strconv.Itoa(int(acc.ID))
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

	res := a.serverCache.Get(ctx, req.GetPhoneNumber())
	if res.Err() != nil {
		logger.WithError(err).Errorln("get otp from redis fail")
		return false, err
	}

	if res.Val() != req.GetOtp() {
		logger.Errorln("wrong otp")
		return false, err
	}
	_, err = a.accountRepo.UpdateStatus(
		ctx,
		req.PhoneNumber,
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
	serverCache *cache.RedisClient,
	accountRepo repositories.AccountRepository,
	authConfig *config.AuthenticationConfig,
) AuthenService {
	return &authenServiceImpl{
		serverCache: serverCache,
		jWTAuth:     jWTAuth,
		accountRepo: accountRepo,
		authConfig:  authConfig,
	}
}
