package server

import (
	"chat-app/config"
	_const "chat-app/internal/const"
	"chat-app/internal/controller"
	repositories2 "chat-app/internal/repositories"
	"chat-app/internal/service"
	"chat-app/internal/webrtc"
	"chat-app/internal/ws"
	"chat-app/pkg/client/cloudinary"
	"chat-app/pkg/client/firebase"
	"chat-app/pkg/client/redis"
	"chat-app/pkg/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"github.com/whatvn/denny"
	dennyHttp "github.com/whatvn/denny/middleware/http"
	"net/http"
	"time"
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
	fb := firebase.GetFirebase(_const.GoogleAccountCertPath)

	chatAppDB := repositories2.InitChatAppDatabase()
	chatMessageDB := repositories2.InitChatMessageDatabase()
	userRepo := repositories2.NewUserRepository(chatAppDB)
	accountRepo := repositories2.NewUserRepository(chatAppDB)
	messageRepo := repositories2.NewMessageRepository(chatMessageDB.DB)
	roomRepo := repositories2.NewRoomRepository(chatAppDB)
	groupRepo := repositories2.NewGroupRepository(chatAppDB)
	fileRepo := repositories2.NewFileRepository(chatAppDB)
	friendRepo := repositories2.NewFriendRepository(chatAppDB)
	// Websockets Setup
	hub := ws.NewHub(redisCli)
	go hub.Run()
	g.GET("/ws/:user_id", func(c *gin.Context) {
		userId := c.Param("user_id")
		ws.ServeWs(c, hub, userId)
	})
	g.GET("/room/create", func(c *gin.Context) {
		c.Redirect(
			http.StatusMovedPermanently,
			fmt.Sprintf("/room/%s/websocket", guuid.New().String()),
		)
	})
	//g.GET("/room/:uuid", func(c *gin.Context) {
	//
	//})
	g.GET("/room/:uuid/websocket", func(c *gin.Context) {
		uuid := c.Param("uuid")
		webrtc.RoomWebsocket(c, uuid)
	})
	g.GET("/room/:uuid/viewer/websocket", func(c *gin.Context) {
		uuid := c.Param("uuid")
		webrtc.RoomViewerWebsocket(c, uuid)
	})

	webrtc.Rooms = make(map[string]*webrtc.CallRoom)
	go dispatchKeyFrames()

	//g.POST("/webrtc/sdp/m/:meetingId/c/:userID/p/:peerId/s/:isSender", func(c *gin.Context) {
	//	isSender, _ := strconv.ParseBool(c.Param("isSender"))
	//	userID := c.Param("userID")
	//	peerID := c.Param("peerId")
	//	webrtc.NewCall(c, isSender, userID, peerID)
	//})

	socketService := service.NewSocketService(hub)
	mailService := service.NewMailService(config.GetAppConfig().Mail, _const.MailTemplatePath)
	fileService := service.NewFileService(cld, fileRepo)
	userService := service.NewUserService(userRepo, fileService, friendRepo)
	friendService := service.NewFriendService(friendRepo, userService)
	roomService := service.NewRoomService(roomRepo)
	groupService := service.NewGroupService(groupRepo, roomRepo, fileService)
	notificationService := service.NewNotificationService(fb, userService, redisCli)
	messageService := service.NewMessageService(
		messageRepo,
		socketService,
		notificationService,
		roomService,
		userService,
		groupService,
	)
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
			roomService,
		),
	)
	apiGroup.Use(middleware.HTTPAuthentication)
	apiGroup.BrpcController(
		controller.NewMessage(
			messageService,
			fileService,
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
	apiGroup.BrpcController(
		controller.NewFriend(
			friendService,
		),
	)
}

func dispatchKeyFrames() {
	for range time.NewTicker(time.Second * 3).C {
		for _, room := range webrtc.Rooms {
			room.Peers.DispatchKeyFrame()
		}
	}
}
