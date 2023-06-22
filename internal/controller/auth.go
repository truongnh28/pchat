package controller

import (
	"chat-app/config"
	"chat-app/helper"
	"chat-app/internal/common"
	"chat-app/internal/service"
	cache "chat-app/pkg/client/redis"
	chat_app "chat-app/proto/chat-app"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/whatvn/denny"
	"net/http"
	"time"
)

type auth struct {
	authService service.AuthService
	config      *config.AuthenticationConfig
	redisCli    cache.IRedisClient
	mailService service.MailService
}

func (a *auth) ResendOtp(
	ctx context.Context,
	request *chat_app.ResendOtpRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode = common.OK
		logger  = denny.GetLogger(ctx)
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()

	resp = new(chat_app.BasicResponse)
	res := a.redisCli.Get(ctx, request.GetEmail())
	if !(res.Err() != nil && res.Err() == redis.Nil) {
		logger.WithError(err).Errorln("otp exist")
		errCode = common.InvalidRequest
		return
	}
	opt, err := helper.GenOtp(config.GetAppConfig().Authentication.SecretKey)
	if err != nil {
		logger.WithError(err).Errorln("gen otp fail")
		errCode = common.SystemError
		return
	}
	err = a.redisCli.Set(ctx, request.GetEmail(), opt, time.Minute*5).Err()
	if err != nil {
		logger.WithError(err).Errorln("gen otp fail")
		errCode = common.SystemError
		return
	}

	err = a.mailService.SendOpt([]string{request.GetEmail()}, opt)
	if err != nil {
		logger.WithError(err).Error("send otp fail")
		errCode = common.SystemError
		return
	}

	return
}

func (a *auth) Logout(
	ctx context.Context,
	request *chat_app.EmptyRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode = common.OK
		logger  = denny.GetLogger(ctx)
		ok      = false
		httpCtx *denny.Context
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
	logger.Infoln("set cookie is empty")
	httpCtx.SetSameSite(http.SameSiteNoneMode)
	httpCtx.SetCookie(
		common.CookieName,
		"", 0, "/", a.config.CookiePath,
		a.config.CookieSecure, false)
	return
}

func (a *auth) VerifyOtp(
	ctx context.Context,
	request *chat_app.VerifyOtpRequest,
) (resp *chat_app.BasicResponse, err error) {
	var (
		errCode = common.OK
		logger  = denny.GetLogger(ctx)
		ok      = false
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.BasicResponse)
	ok, err = a.authService.VerifyOpt(ctx, request)
	if err != nil {
		errCode = common.SystemError
		logger.WithError(err).Errorln("verify otp fail err")
		return
	}
	if !ok {
		errCode = common.ValidationError
		logger.Infoln("validate otp fail")
		return
	}
	return
}

func (a *auth) Login(
	ctx context.Context,
	request *chat_app.LoginRequest,
) (resp *chat_app.LoginResponse, err error) {
	var (
		errCode = common.OK
		logger  = denny.GetLogger(ctx)
		httpCtx *denny.Context
		ok      = false
	)
	defer func() {
		buildResponse(errCode, resp)
		err = nil
	}()
	resp = new(chat_app.LoginResponse)
	httpCtx, ok = ctx.(*denny.Context)
	if !ok {
		errCode = common.SystemError
		logger.WithError(errors.New("get httpCtx fail"))
		return
	}
	resp, err = a.authService.Login(ctx, request)

	if err != nil {
		errCode = common.InvalidRequest
		logger.WithError(err).Errorln("authen service fail")
		return
	}

	httpCtx.SetCookie(
		common.CookieName,
		resp.AccessToken,
		int(a.config.ExpiredTime), "/", "",
		httpCtx.Request.TLS != nil, false)
	resp.AccessToken = ""
	return
}

func NewAuth(
	authService service.AuthService,
	redisCli cache.IRedisClient,
	mailService service.MailService,
	authConfig *config.AuthenticationConfig,
) chat_app.AuthServer {
	return &auth{
		authService: authService,
		config:      authConfig,
		redisCli:    redisCli,
		mailService: mailService,
	}
}
