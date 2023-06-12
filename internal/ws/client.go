package ws

import (
	"chat-app/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/whatvn/denny"
	"net/http"
	"time"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second
	// Max time till next pong from peer
	pongWait = 60 * time.Second
	// Send ping interval, must be less than pong wait time
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var endline = []byte{'\n'}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// get socket context
var ctx = context.Background()

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id          string
	hub         *Hub
	wsConn      *websocket.Conn
	sendChannel chan []byte
}

// ServeWs handles websockets requests from clients requests.
func ServeWs(ctx *gin.Context, hub *Hub, userId string) interface{} {
	var (
		logger   = denny.GetLogger(ctx)
		username = userId
	)
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.WithError(err).Errorln("start ws fail")
		return nil
	}

	client := newClient(conn, hub, username)
	go client.writePump()
	go client.readPump()
	// sendChannel this user to register chanel
	client.hub.RegisterChannel <- client
	logger.Infof("ServeWs of user %s start!!!\n", username)
	return nil
}

func newClient(conn *websocket.Conn, hub *Hub, username string) *Client {
	return &Client{
		id:          username,
		hub:         hub,
		wsConn:      conn,
		sendChannel: make(chan []byte, 256),
	}
}
func unregisterAndCloseConnection(c *Client) {
	c.hub.UnregisterChannel <- c
	_ = c.wsConn.Close()
}

func setSocketPayloadConfig(c *Client) {
	// set the maximum message size limit that the client can receive from the server.
	c.wsConn.SetReadLimit(maxMessageSize)

	// Once the latest read time is set, the client will wait for the next message from the server within this interval.
	// If no messages are received within this time, the connection will be closed.
	_ = c.wsConn.SetReadDeadline(time.Now().Add(pongWait))

	// It is used to determine the maximum time the client should wait before receiving the next "pong" message
	// from the server. If this time is exceeded, the connection will be closed.
	c.wsConn.SetPongHandler(func(appData string) error {
		_ = c.wsConn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

func (client *Client) disconnect() {
	client.hub.UnregisterChannel <- client
	close(client.sendChannel)
	_ = client.wsConn.Close()
}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = client.wsConn.Close()
	}()
	for {
		select {
		case payload, ok := <-client.sendChannel:
			//requestBody := new(bytes.Buffer)
			//_ = json.NewEncoder(requestBody).Encode(payload)
			//finalPayload := requestBody.Bytes()
			//fmt.Println(string(finalPayload))
			// set time write message, if timeout -> close connect
			_ = client.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = client.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := client.wsConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, _ = w.Write(payload)
			n := len(client.sendChannel)

			for i := 0; i < n; i++ {
				_, _ = w.Write(endline)
				_, _ = w.Write(<-client.sendChannel)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = client.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) readPump() {
	logger := denny.GetLogger(ctx)
	// unregister and close connection
	defer func() {
		client.disconnect()
	}()

	// Setting up the Payload configuration
	setSocketPayloadConfig(client)

	for {
		_, payload, err := client.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				logger.WithError(err).Errorf("read message from web socket fail")
				break
			}
		}

		//Getting the proper payload to sendChannel the client
		client.handleNewMessage(payload)
	}
}

func (client *Client) handleNewMessage(msg []byte) {
	var message ReceivedMessage
	logger := denny.GetLogger(ctx)
	if err := json.Unmarshal(msg, &message); err != nil {
		logger.WithError(err).Error("unmarshal JSON message fail")
	}
	switch message.Event {
	case JoinUser:
		client.handleJoinGroupMessage(message)
	case NewMessage:
		client.handleEmitMessage(message)
	case DeleteMessage:

	case LeaveUser:
		client.handleLeaveGroup(message)
	case JoinGroup:
	case LeaveGroup:
	case StartTyping:
	case StopTyping:

	}
}

type chatListResponseStruct struct {
	Type     string      `json:"type"`
	ChatList interface{} `json:"chat_list"`
}

func (client *Client) handleEmitMessage(message ReceivedMessage) {
	fmt.Println("handleEmitMessage")
	logger := denny.GetLogger(ctx).WithField("socket payload event", client.id)

	msg := (message.Payload).(map[string]interface{})["message"].(string)
	senderId := (message.Payload).(map[string]interface{})["senderId"].(string)
	recipientId := (message.Payload).(map[string]interface{})["recipientId"].(string)
	if !(msg != "" && senderId != "" && recipientId != "") {
		logger.WithError(errors.New("get data from payload fail")).
			Error("get data from payload fail")
		return
	}
	messagePacket := domain.ChatMessage{
		SenderID:    senderId,
		RecipientID: recipientId,
		Message:     msg,
		Time:        time.Now(),
	}
	err := client.hub.MessageService.StoreNewChatMessages(ctx, &messagePacket)
	if err != nil {
		fmt.Println(err)
		logger.WithError(err).Errorln("store message fail")
	}
	e, err := sonic.Marshal(SocketMessage{
		Event:   NewMessage,
		Payload: messagePacket,
	})
	if err != nil {
		logger.WithError(err).Errorf("marshal msg to byte array fail")
		return
	}
	client.hub.EmitToSpecificClient(e, recipientId)
}

func (client *Client) handleJoinGroupMessage(message ReceivedMessage) {
	logger := denny.GetLogger(ctx)
	userId := (message.Payload).(string)
	userDetail, err := client.hub.UserService.GetByUserId(userId)
	if err != nil {
		logger.WithError(err).Errorf("get user by id fail")
		return
	}

	msg, err := sonic.Marshal(SocketMessage{
		Event: JoinUser,
		Payload: domain.UserDetail{
			ID:       userId,
			Username: userDetail.Username,
			Online:   userDetail.Online,
		},
	})

	if err != nil {
		logger.WithError(err).Errorf("marshal message to byte array fail")
		return
	}

	client.hub.BroadcastSocketEventToAllClientExceptMe(msg, userId)

	userOnl, err := client.hub.UserService.GetAllOnlineUsers(userId)
	if err != nil {
		logger.WithError(err).Errorf("marshal message to byte array fail")
		return
	}

	allOnlineUsers := SocketMessage{
		Event: "chatlist-response",
		Payload: chatListResponseStruct{
			Type:     "my-chat-list",
			ChatList: userOnl,
		},
	}
	msg, err = sonic.Marshal(allOnlineUsers)
	if err != nil {
		logger.WithError(err).Errorf("marshal message to byte array fail")
		return
	}
	client.hub.EmitToSpecificClient(msg, userId)
}

func (client *Client) handleLeaveGroup(message ReceivedMessage) {
	userId := (message.Payload).(string)
	logger := denny.GetLogger(ctx)
	userDetail, err := client.hub.UserService.GetByUserId(userId)
	if err != nil {
		logger.WithError(err).Errorf("get user by id fail")
		return
	}
	_ = client.hub.UserService.UpdateUserOnlineStatusByUserID(userId, false)
	msg, err := sonic.Marshal(SocketMessage{
		Event: "chatlist-response",
		Payload: chatListResponseStruct{
			Type: "user-disconnect",
			ChatList: domain.UserDetail{
				ID:       userId,
				Username: userDetail.Username,
				Online:   false,
			},
		},
	})
	if err != nil {
		logger.WithError(err).Errorf("marshal message to byte array fail")
		return
	}
	client.hub.BroadcastSocketEventToAllClient(msg)
}
