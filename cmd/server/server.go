package server

import (
	"chat-app/config"
	_const "chat-app/internal/const"
	"chat-app/internal/controller"
	"chat-app/internal/service"
	"chat-app/internal/ws"
	"chat-app/pkg/repositories"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/whatvn/denny"
	dennyHttp "github.com/whatvn/denny/middleware/http"
)

type HTTPServer struct {
}

func (h *HTTPServer) Init() {
	initAuthentication()
}

func (_ *HTTPServer) Run() error {
	serverConfig := config.GetAppConfig().Server

	server := denny.NewServer(serverConfig.Debug)
	server.WithMiddleware(gin.Recovery(), dennyHttp.Logger())
	server.RedirectTrailingSlash = false
	setupHandler(server)
	server.Info("starting http server...")
	return server.GraceFulStart(fmt.Sprintf(":%s", serverConfig.Port))
}

func initAuthentication() {
	//TODO: handle auth
}

func setupHandler(s *denny.Denny) {
	g := s.Engine

	authGroup := s.NewGroup("/api")
	authGroup.WithCors()

	apiGroup := s.NewGroup("/api")
	apiGroup.WithCors()
	wsGroup := s.NewGroup("/ws")
	wsGroup.WithCors()

	chatAppDB := repositories.InitChatAppDatabase()
	chatMessageDB := repositories.InitChatMessageDatabase()
	userRepo := repositories.NewUserRepository(chatAppDB)
	userService := service.NewUserService(userRepo)
	messageRepo := repositories.NewMessageRepository(chatMessageDB.DB)
	messageService := service.NewMessageService(messageRepo)
	accountRepo := repositories.NewAccountRepository(chatAppDB)
	accountService := service.NewAccountService(accountRepo)
	mailService := service.NewMailService(config.GetAppConfig().Mail, _const.MailTemplatePath)

	apiGroup.BrpcController(
		controller.NewMessage(
			messageService,
		),
	)

	apiGroup.BrpcController(
		controller.NewUser(
			userService,
		),
	)
	apiGroup.BrpcController(
		controller.NewAccount(
			accountService,
			mailService,
		),
	)
	// Websockets Setup
	hub := ws.NewHub(userService, messageService)
	go hub.Run()
	g.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(c, hub)
	})
}
