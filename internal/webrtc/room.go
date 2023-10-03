package webrtc

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/whatvn/denny"
	"log"
	"sync"
)

func RoomConn(ctx context.Context, c *websocket.Conn, p *Peers) {
	var (
		logger = denny.GetLogger(ctx)
		config webrtc.Configuration
	)
	//if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
	//	config = turnConfig
	//}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		logger.WithError(err).Errorln("new peer connection err: ", err)
		return
	}
	defer peerConnection.Close()

	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			logger.WithError(err).Errorln("add transceiver from kind err: ", err)
			return
		}
	}

	newPeer := PeerConnectionState{
		PeerConnection: peerConnection,
		Websocket: &ThreadSafeWriter{
			Conn:  c,
			Mutex: sync.Mutex{},
		}}

	// Add our new PeerConnection to global list
	p.ListLock.Lock()
	p.Connections = append(p.Connections, newPeer)
	p.ListLock.Unlock()

	logger.Infoln("connection: ", p.Connections)

	// Trickle ICE. Emit server candidate to client
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidateString, err := sonic.Marshal(i.ToJSON())
		if err != nil {
			logger.WithError(err).Errorln("marshal err: ", err)
			return
		}

		if writeErr := newPeer.Websocket.WriteJSON(&websocketMessage{
			Event: Candidate,
			Data:  string(candidateString),
		}); writeErr != nil {
			logger.WithError(writeErr).Errorln("write json err: ", err)
		}
	})

	// If PeerConnection is closed remove it from global list
	peerConnection.OnConnectionStateChange(func(pp webrtc.PeerConnectionState) {
		switch pp {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				logger.WithError(err)
			}
		case webrtc.PeerConnectionStateClosed:
			p.SignalPeerConnections()
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		// Create a track to fan out our incoming video to all peers
		trackLocal := p.AddTrack(t)
		if trackLocal == nil {
			return
		}
		defer p.RemoveTrack(trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := t.Read(buf)
			if err != nil {
				return
			}

			if _, err = trackLocal.Write(buf[:i]); err != nil {
				return
			}
		}
	})

	p.SignalPeerConnections()
	message := &websocketMessage{}
	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		} else if err := sonic.Unmarshal(raw, &message); err != nil {
			logger.WithError(err).Errorln("unmarshal err: ", err)
			return
		}

		switch message.Event {
		case Candidate:
			candidate := webrtc.ICECandidateInit{}
			if err := sonic.Unmarshal([]byte(message.Data), &candidate); err != nil {
				logger.WithError(err).Errorln("unmarshal err: ", err)
				return
			}

			if err := peerConnection.AddICECandidate(candidate); err != nil {
				logger.WithError(err).Errorln("add ICE candidate err: ", err)
				return
			}
		case Answer:
			answer := webrtc.SessionDescription{}
			if err := sonic.Unmarshal([]byte(message.Data), &answer); err != nil {
				logger.WithError(err).Errorln("unmarshal err: ", err)
				return
			}

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				logger.WithError(err).Errorln("set remote description err: ", err)
				return
			}
		}
	}
}
