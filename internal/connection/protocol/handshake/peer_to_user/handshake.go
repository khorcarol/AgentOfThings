package peer_to_user

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/khorcarol/AgentOfThings/internal/api"
)

const (
	HandshakeTimeout    = 10 * time.Second
	maxInterests        = 10 // Arbitrarily chosen for now, require discussion with team to properly choose.
	HandshakeProtocolID = "agentofthings/peer_to_user_handshake/0.0.1"
)

// InitiateHandshake dials a remote peer, sends our interests (as a JSON HandshakeMessage),
// half-closes the write side, then waits for a response.
// It returns the remote peer’s handshake message on success.
func InitiateHandshake(host host.Host, ctx context.Context, remote peer.ID, usr api.User) (*api.User, error) {
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

	if len(usr.Interests) > maxInterests {
		usr.Interests = usr.Interests[:maxInterests]
	}

	// Send our interests initially.
	if err := encodeToStream(stream, usr); err != nil {
		stream.Reset()
		return nil, fmt.Errorf("handshake: failed to encode our handshake message: %w", err)
	}

	// Half-close the write side as no more writes from us needed.
	if halfCloser, ok := stream.(interface{ CloseWrite() error }); ok {
		if err := halfCloser.CloseWrite(); err != nil {
			log.Printf("handshake: failed to half-close write: %v", err)
		}
	}

	log.Printf("handshake: sent our handshake message: %+v", usr)

	var remoteMessage api.User
	if err := decodeFromStream(stream, &remoteMessage); err != nil {
		stream.Reset()
		return nil, fmt.Errorf("handshake: failed to decode remote handshake message: %w", err)
	}
	log.Printf("handshake: received handshake message: %+v", &remoteMessage)
	return &remoteMessage, nil
}

// handshakeHandler is invoked when a remote peer connects.
// It decodes the remote’s handshake message and responds with our interests.
func HandshakeHandler(stream network.Stream, callback func(*api.User, peer.ID), toSend *api.User) {
	defer stream.Close()

	if err := stream.SetDeadline(time.Now().Add(HandshakeTimeout)); err != nil {
		log.Printf("handshake: error setting deadline: %v", err)
		stream.Reset()
		return
	}

	var remoteMessage api.User
	if err := decodeFromStream(stream, &remoteMessage); err != nil {
		log.Printf("handshake: failed to decode remote handshake message: %v", err)
		stream.Reset()
		return
	}
	log.Printf("handshake: received handshake message from peer: %+v", &remoteMessage)

	ourMessage := *toSend
	if len(ourMessage.Interests) > maxInterests {
		ourMessage.Interests = ourMessage.Interests[:maxInterests]
	}

	if err := encodeToStream(stream, ourMessage); err != nil {
		log.Printf("handshake: failed to encode our handshake message: %v", err)
		stream.Reset()
		return
	}

	// Half-close our write side.
	if halfCloser, ok := stream.(interface{ CloseWrite() error }); ok {
		if err := halfCloser.CloseWrite(); err != nil {
			log.Printf("handshake: failed to half-close write: %v", err)
		}
	}
	log.Printf("handshake: sent our handshake message: %+v", ourMessage)
	callback(&remoteMessage, stream.Conn().RemotePeer())
}

// encodeToStream marshals the given message to JSON with a length prefix and writes it to w.
func encodeToStream(writer io.Writer, message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	length := uint32(len(data))
	if err := binary.Write(writer, binary.LittleEndian, length); err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

// decodeFromStream reads a length-prefixed JSON message from r into the provided message.
func decodeFromStream(reader io.Reader, message any) error {
	var length uint32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return err
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(reader, data); err != nil {
		return err
	}
	return json.Unmarshal(data, message)
}

func ShouldInitiate(localPeerID, remotePeerID peer.ID) bool {
	return localPeerID.String() < remotePeerID.String()
}
