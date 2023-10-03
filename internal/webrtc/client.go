package webrtc

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/whatvn/denny"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RoomWebsocket(ctx *gin.Context, uuid string) {
	var (
		logger = denny.GetLogger(ctx)
	)
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.WithError(err).Errorln("start ws fail")
		return
	}

	_, room := createOrGetRoom(uuid)
	RoomConn(ctx, conn, room.Peers)

	return
}

func createOrGetRoom(uuid string) (string, *CallRoom) {
	RoomsLock.Lock()
	defer RoomsLock.Unlock()

	h := sha256.New()
	h.Write([]byte(uuid))

	p := &Peers{}
	p.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
	room := &CallRoom{
		Peers: p,
	}

	Rooms[uuid] = room

	return uuid, room
}

func RoomViewerWebsocket(ctx *gin.Context, uuid string) {
	var (
		logger = denny.GetLogger(ctx)
	)
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.WithError(err).Errorln("start ws fail")
		return
	}

	RoomsLock.Lock()
	if peer, ok := Rooms[uuid]; ok {
		RoomsLock.Unlock()
		roomViewerConn(conn, peer.Peers)
		return
	}
	RoomsLock.Unlock()
}

func roomViewerConn(c *websocket.Conn, p *Peers) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer c.Close()

	for {
		select {
		case <-ticker.C:
			w, err := c.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(fmt.Sprintf("%d", len(p.Connections))))
		}
	}
}
