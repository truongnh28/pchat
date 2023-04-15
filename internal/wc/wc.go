package wc

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

var (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = int64(512)
)

// NewSocketUser this function will create a new user/client
func NewSocketUser(hub *Hub, conn *websocket.Conn, username string) {
	client := &Client{
		hub:         hub,
		wsConn:      conn,
		sendChannel: make(chan SocketEventStruct),
		userId:      username,
	}

	go client.writePump()
	go client.readPump()
	// sendChannel this user to register chanel
	client.hub.RegisterChannel <- client
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

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.wsConn.Close()
	}()
	for {
		select {
		case payload, ok := <-c.sendChannel:
			requestBody := new(bytes.Buffer)
			_ = json.NewEncoder(requestBody).Encode(payload)
			finalPayload := requestBody.Bytes()
			// set time write message, if timeout -> close connect
			_ = c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.wsConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(finalPayload)
			n := len(c.sendChannel)
			for i := 0; i < n; i++ {
				_ = json.NewEncoder(requestBody).Encode(<-c.sendChannel)
				_, _ = w.Write(requestBody.Bytes())
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	var socketEventPayload SocketEventStruct
	// unregister and close connection
	defer unregisterAndCloseConnection(c)

	// Setting up the Payload configuration
	setSocketPayloadConfig(c)

	for {
		_, payload, err := c.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// TODO: handle log
				break
			}
		}
		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoderErr := decoder.Decode(&socketEventPayload)
		if decoderErr != nil {
			// TODO: handle log
			break
		}

		//Getting the proper payload to sendChannel the client
		handleSocketPayloadEvent(c, &socketEventPayload)
	}
}
