package send_message

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection/protocol/transport"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

const (
	ProtocolTimeout = 10 * time.Second
	ProtocolID      = "agentofthings/send_message/0.0.1"
)

func SendMessages(host host.Host, ctx context.Context, remote peer.ID, toSend []api.Message) error {
	stream, err := host.NewStream(ctx, remote, protocol.ID(ProtocolID))
	if err != nil {
		return fmt.Errorf("send_message: failed to open stream to peer %s: %w", remote, err)
	}
	defer stream.Close()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(ProtocolTimeout)
	}
	if err := stream.SetDeadline(deadline); err != nil {
		return fmt.Errorf("send_message: setting deadline failed: %v", err)
	}

	if err := transport.EncodeToStream(stream, toSend); err != nil {
		stream.Reset()
		return fmt.Errorf("send_message: failed to encode our handshake message: %w", err)
	}

	if halfCloser, ok := stream.(interface{ CloseWrite() error }); ok {
		if err := halfCloser.CloseWrite(); err != nil {
			return fmt.Errorf("send_message: failed to half-close write: %v", err)
		}
	}

	return nil
}

func SendMessageHandler(stream network.Stream, callback func(*[]api.Message, peer.ID)) {
	defer stream.Close()

	if err := stream.SetDeadline(time.Now().Add(ProtocolTimeout)); err != nil {
		log.Printf("send_message handler: error setting deadline: %v", err)
		stream.Reset()
		return
	}

	var remoteMessage []api.Message
	if err := transport.DecodeFromStream(stream, &remoteMessage); err != nil {
		log.Printf("send_message handler: failed to decode remote handshake message: %v", err)
		stream.Reset()
		return
	}

	if halfCloser, ok := stream.(interface{ CloseWrite() error }); ok {
		if err := halfCloser.CloseWrite(); err != nil {
			log.Printf("send_message handler: failed to half-close write: %v", err)
		}
	}
	callback(&remoteMessage, stream.Conn().RemotePeer())
}
