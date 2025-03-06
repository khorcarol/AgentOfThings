package user_to_friend

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection/protocol/transport"
)

const FriendRequestProtocolID = "agentofthings/friend_request/0.0.1"

// SendFriendData opens a stream to send friend-related data (either a request or a response),
// writes the data, half-closes the stream, and returns any error received.
func SendFriendData(host host.Host, ctx context.Context, remote peer.ID, friendData api.FriendRequest) error {
	stream, err := host.NewStream(ctx, remote, protocol.ID(FriendRequestProtocolID))
	if err != nil {
		return fmt.Errorf("friend data: failed to open stream to peer %s: %w", remote, err)
	}
	defer stream.Close()

	if err := transport.EncodeToStream(stream, friendData); err != nil {
		stream.Reset()
		return fmt.Errorf("friend data: failed to encode friend data: %w", err)
	}
	if halfCloser, ok := stream.(interface{ CloseWrite() error }); ok {
		if err := halfCloser.CloseWrite(); err != nil {
			log.Printf("friend data: failed to half-close write side: %v", err)
		}
	}
	log.Printf("friend data: sent data: %+v", friendData)
	return nil
}

// FriendRequestHandler handles incoming friend data over the friend protocol.
// It decodes the data into an api.Friend and then passes it along via callback.
// The [callback] function should decide if this friend is to be displayed or stored.
func FriendRequestHandler(stream network.Stream, callback func(*api.FriendRequest, peer.ID)) {
	defer stream.Close()

	var incomingFriend api.FriendRequest
	if err := transport.DecodeFromStream(stream, &incomingFriend); err != nil {
		log.Printf("friend data handler: failed to decode friend data: %v", err)
		stream.Reset()
		return
	}
	log.Printf("friend data handler: received friend data: %+v", incomingFriend)

	callback(&incomingFriend, stream.Conn().RemotePeer())
}
