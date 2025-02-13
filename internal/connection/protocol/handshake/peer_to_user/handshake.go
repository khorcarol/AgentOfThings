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
	"github.com/khorcarol/AgentOfThings/internal/api/interests"
)

const (
	HandshakeTimeout    = 10 * time.Second
	maxInterests        = 10 // Arbitrarily chosen for now, require discussion with team to properly choose.
	HandshakeProtocolID = "agentofthings/peer_to_user_handshake/0.0.1"
)

// HandshakeMessage represents the message exchanged during a handshake.
// It contains a slice of Interests (from our internal API package).
type HandshakeMessage struct {
	Interests []*api.Interest `json:"interests"`
}

// HandshakeService encapsulates the host for the handshake protocol.
type HandshakeService struct {
	Host host.Host
}

// NewHandshakeService creates a new handshake service and registers the stream handler.
func NewHandshakeService(host host.Host) *HandshakeService {
	service := &HandshakeService{Host: host}
	host.SetStreamHandler(protocol.ID(HandshakeProtocolID), service.handshakeHandler)
	return service
}

// InitiateHandshake dials a remote peer, sends our interests (as a JSON HandshakeMessage),
// half-closes the write side, then waits for a response.
// It returns the remote peer’s handshake message on success.
func (service *HandshakeService) InitiateHandshake(ctx context.Context, remote peer.ID) (*HandshakeMessage, error) {
	stream, err := service.Host.NewStream(ctx, remote, protocol.ID(HandshakeProtocolID))
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

	ourMessage := getLocalHandshakeMessage()
	if len(ourMessage.Interests) > maxInterests {
		ourMessage.Interests = ourMessage.Interests[:maxInterests]
	}

	// Send our interests initially.
	if err := encodeToStream(stream, ourMessage); err != nil {
		stream.Reset()
		return nil, fmt.Errorf("handshake: failed to encode our handshake message: %w", err)
	}

	// Half-close the write side as no more writes from us needed.
	if halfCloser, ok := stream.(interface{ CloseWrite() error }); ok {
		if err := halfCloser.CloseWrite(); err != nil {
			log.Printf("handshake: failed to half-close write: %v", err)
		}
	}

	log.Printf("handshake: sent our handshake message: %+v", ourMessage)

	var remoteMessage HandshakeMessage
	if err := decodeFromStream(stream, &remoteMessage); err != nil {
		stream.Reset()
		return nil, fmt.Errorf("handshake: failed to decode remote handshake message: %w", err)
	}
	log.Printf("handshake: received handshake message: %+v", &remoteMessage)
	return &remoteMessage, nil
}

// handshakeHandler is invoked when a remote peer connects.
// It decodes the remote’s handshake message and responds with our interests.
func (service *HandshakeService) handshakeHandler(stream network.Stream) {
	defer stream.Close()

	if err := stream.SetDeadline(time.Now().Add(HandshakeTimeout)); err != nil {
		log.Printf("handshake: error setting deadline: %v", err)
		stream.Reset()
		return
	}

	// TODO: write remoteMessage to some channel to get the data to the backend
	var remoteMessage HandshakeMessage
	if err := decodeFromStream(stream, &remoteMessage); err != nil {
		log.Printf("handshake: failed to decode remote handshake message: %v", err)
		stream.Reset()
		return
	}
	log.Printf("handshake: received handshake message from peer: %+v", &remoteMessage)

	ourMessage := getLocalHandshakeMessage()
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
}

// encodeToStream marshals the given message to JSON with a length prefix and writes it to w.
func encodeToStream(writer io.Writer, message interface{}) error {
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
func decodeFromStream(reader io.Reader, message interface{}) error {
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

// Just a silly stub for now.
func getLocalHandshakeMessage() *HandshakeMessage {
	return &HandshakeMessage{
		Interests: []*api.Interest{
			{
				Category:    interests.Sport,
				Description: "RTS, CTS, ACK Bryan grrrrrrrr",
			},
			{
				Category:    interests.Music,
				Description: "Jazz and classical tunes.",
			},
		},
	}
}
