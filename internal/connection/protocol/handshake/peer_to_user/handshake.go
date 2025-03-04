package peer_to_user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection/protocol/transport"
	"github.com/khorcarol/AgentOfThings/internal/personal"
)

const (
	HandshakeTimeout    = 10 * time.Second
	maxInterests        = 10 // Arbitrarily chosen for now, require discussion with team to properly choose.
	HandshakeProtocolID = "agentofthings/peer_to_user_handshake/0.0.1"
)

// InitiateHandshake dials a remote peer, sends our interests (as a JSON HandshakeMessage),
// half-closes the write side, then waits for a response.
// It returns the remote peer’s handshake message on success.
func InitiateHandshake(host host.Host, ctx context.Context, remote peer.ID) (*api.User, error) {
	stream, err := host.NewStream(ctx, remote, protocol.ID(HandshakeProtocolID))
	if err != nil {
		return nil, fmt.Errorf("handshake: failed to open stream to peer %s: %w", remote, err)
	}
	defer stream.Close()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(HandshakeTimeout)
	}
	if err := stream.SetDeadline(deadline); err != nil {
		log.Printf("handshake: setting deadline failed: %v", err)
	}

	friend := personal.GetSelf()

	usr := friend.User
	if len(usr.Interests) > maxInterests {
		usr.Interests = usr.Interests[:maxInterests]
	}

	if err := transport.EncodeToStream(stream, usr); err != nil {
		stream.Reset()
		return nil, fmt.Errorf("handshake: failed to encode our handshake message: %w", err)
	}

	if halfCloser, ok := stream.(interface{ CloseWrite() error }); ok {
		if err := halfCloser.CloseWrite(); err != nil {
			log.Printf("handshake: failed to half-close write: %v", err)
		}
	}

	log.Printf("handshake: sent our handshake message: %+v", usr)

	var remoteMessage api.User
	if err := transport.DecodeFromStream(stream, &remoteMessage); err != nil {
		stream.Reset()
		return nil, fmt.Errorf("handshake: failed to decode remote handshake message: %w", err)
	}
	log.Printf("handshake: received handshake message: %+v", &remoteMessage)
	return &remoteMessage, nil
}

// handshakeHandler is invoked when a remote peer connects.
// It decodes the remote’s handshake message and responds with our interests.
func HandshakeHandler(stream network.Stream, callback func(*api.User, peer.ID)) {
	defer stream.Close()

	if err := stream.SetDeadline(time.Now().Add(HandshakeTimeout)); err != nil {
		log.Printf("handshake: error setting deadline: %v", err)
		stream.Reset()
		return
	}

	var remoteMessage api.User
	if err := transport.DecodeFromStream(stream, &remoteMessage); err != nil {
		log.Printf("handshake: failed to decode remote handshake message: %v", err)
		stream.Reset()
		return
	}
	log.Printf("handshake: received handshake message from peer: %+v", &remoteMessage)

	toSend := personal.GetSelf()

	ourMessage := toSend.User
	if len(ourMessage.Interests) > maxInterests {
		ourMessage.Interests = ourMessage.Interests[:maxInterests]
	}

	if err := transport.EncodeToStream(stream, ourMessage); err != nil {
		log.Printf("handshake: failed to encode our handshake message: %v", err)
		stream.Reset()
		return
	}

	if halfCloser, ok := stream.(interface{ CloseWrite() error }); ok {
		if err := halfCloser.CloseWrite(); err != nil {
			log.Printf("handshake: failed to half-close write: %v", err)
		}
	}
	log.Printf("handshake: sent our handshake message: %+v", ourMessage)
	callback(&remoteMessage, stream.Conn().RemotePeer())
}
