package identify

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/khorcarol/AgentOfThings/internal/connection/protocol/transport"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

const (
	HandshakeTimeout    = 10 * time.Second
	HandshakeProtocolID = "agentofthings/identify_handshake/0.0.1"
)

// Determines whether a peer is a user or a hub.
// Returns true for user, false for hub
func Identify(host host.Host, ctx context.Context, remote peer.ID) (bool, error) {
	stream, err := host.NewStream(ctx, remote, protocol.ID(HandshakeProtocolID))
	if err != nil {
		return false, fmt.Errorf("handshake: failed to open stream to peer %s: %w", remote, err)
	}
	defer stream.Close()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(HandshakeTimeout)
	}
	if err := stream.SetDeadline(deadline); err != nil {
		log.Printf("handshake: setting deadline failed: %v", err)
	}
	// Close our write end because we only want to receive data
	if halfCloser, ok := stream.(interface{ CloseWrite() error }); ok {
		if err := halfCloser.CloseWrite(); err != nil {
			log.Printf("handshake: failed to half-close write: %v", err)
		}
	}

	var response bool
	if err := transport.DecodeFromStream(stream, &response); err != nil {
		stream.Reset()
		return false, fmt.Errorf("handshake: failed to decode remote handshake message: %w", err)
	}
	return response, nil
}

func IdentifyHandler(stream network.Stream, isUser bool, callback func(peer.ID)) {
	defer stream.Close()

	if err := stream.SetDeadline(time.Now().Add(HandshakeTimeout)); err != nil {
		log.Printf("handshake: error setting deadline: %v", err)
		stream.Reset()
		return
	}

	if err := transport.EncodeToStream(stream, isUser); err != nil {
		log.Printf("handshake: failed to encode our handshake message: %v", err)
		stream.Reset()
		return
	}

	if halfCloser, ok := stream.(interface{ CloseWrite() error }); ok {
		if err := halfCloser.CloseWrite(); err != nil {
			log.Printf("handshake: failed to half-close write: %v", err)
		}
	}
	callback(stream.Conn().RemotePeer())
}
