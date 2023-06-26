package server

import (
	"chat-app/config"
	_const "chat-app/internal/const"
	"chat-app/internal/controller"
	"chat-app/internal/service"
	"chat-app/internal/ws"
	"chat-app/pkg/client/cloudinary"
	"chat-app/pkg/client/redis"
	"chat-app/pkg/middleware"
	"chat-app/pkg/repositories"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/whatvn/denny"
	dennyHttp "github.com/whatvn/denny/middleware/http"
)

type HTTPServer struct {
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

func setupHandler(s *denny.Denny) {
	g := s.Engine

	apiGroup := s.NewGroup("/api")
	apiGroup.WithCors()
	wsGroup := s.NewGroup("/ws")
	wsGroup.WithCors()

	redisCli := redis.GetRedisClient(config.GetAppConfig().Redis)
	cld := cloudinary.GetCloudinaryAPI(config.GetAppConfig().Cloudinary)

	chatAppDB := repositories.InitChatAppDatabase()
	chatMessageDB := repositories.InitChatMessageDatabase()
	userRepo := repositories.NewUserRepository(chatAppDB)
	accountRepo := repositories.NewUserRepository(chatAppDB)
	messageRepo := repositories.NewMessageRepository(chatMessageDB.DB)
	roomRepo := repositories.NewRoomRepository(chatAppDB)
	groupRepo := repositories.NewGroupRepository(chatAppDB)
	fileRepo := repositories.NewFileRepository(chatAppDB)
	// Websockets Setup
	hub := ws.NewHub(redisCli)
	go hub.Run()
	g.GET("/ws/:user_id", func(c *gin.Context) {
		userId := c.Param("user_id")
		ws.ServeWs(c, hub, userId)
	})

	socketService := service.NewSocketService(hub)
	messageService := service.NewMessageService(messageRepo, socketService)
	mailService := service.NewMailService(config.GetAppConfig().Mail, _const.MailTemplatePath)
	fileService := service.NewFileService(cld, fileRepo)
	userService := service.NewUserService(userRepo, fileService)
	roomService := service.NewRoomService(roomRepo)
	groupService := service.NewGroupService(groupRepo, roomRepo, fileService)
	authService := service.NewAuthenService(
		service.GetJWTInstance(),
		redisCli,
		accountRepo,
		config.GetAppConfig().Authentication,
	)
	apiGroup.BrpcController(
		controller.NewAuth(
			authService,
			userService,
			redisCli,
			mailService,
			config.GetAppConfig().Authentication,
		),
	)
	apiGroup.Use(middleware.HTTPAuthentication)
	apiGroup.BrpcController(
		controller.NewMessage(
			messageService,
			fileService,
			socketService,
			userService,
			roomService,
		),
	)

	apiGroup.BrpcController(
		controller.NewUser(
			userService,
		),
	)

	apiGroup.BrpcController(
		controller.NewGroup(
			userService,
			groupService,
			roomService,
			fileService,
		),
	)
}
