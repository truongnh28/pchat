package server

import (
	"chat-app/config"
	"chat-app/internal/controller"
	"chat-app/internal/service"
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
	authGroup := s.NewGroup("/api")
	authGroup.WithCors()

	apiGroup := s.NewGroup("/api")
	apiGroup.WithCors()
	chatAppDB := repositories.InitChatAppDatabase()
	chatMessageDB := repositories.InitChatMessageDatabase()
	userRepo := repositories.NewUserRepository(chatAppDB)
	messageRepo := repositories.NewMessageRepository(chatMessageDB)
	messageService := service.NewMessageService(messageRepo)
	apiGroup.BrpcController(
		controller.NewMessage(
			messageService,
		),
	)
	apiGroup.BrpcController(
		controller.NewMessage(
			messageService,
		),
	)
	apiGroup.BrpcController(
		controller.NewUser(
			userRepo,
		),
	)
}
