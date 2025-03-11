package connection

import (
	"context"
	"sync"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection/discovery"
	"github.com/khorcarol/AgentOfThings/internal/connection/protocol/handshake/identify"
	"github.com/khorcarol/AgentOfThings/internal/connection/protocol/send_message"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type ConnectionManager struct {
	mDNSservice    mdns.Service
	peerAddrChan   <-chan peer.AddrInfo
	host           host.Host
	peersMutex     sync.Mutex
	connectedPeers map[peer.ID]struct{}
	self           api.Hub

	NewMessages chan api.Message
}

func (cmgr *ConnectionManager) peerDisconnectWrapper() func(network.Network, network.Conn) {
	return func(n network.Network, c network.Conn) {
		cmgr.peersMutex.Lock()
		defer cmgr.peersMutex.Unlock()
		delete(cmgr.connectedPeers, c.RemotePeer())
	}
}

func (cmgr *ConnectionManager) peerConnectWrapper() func(network.Network, network.Conn) {
	return func(n network.Network, c network.Conn) {
		cmgr.peersMutex.Lock()
		defer cmgr.peersMutex.Unlock()
		cmgr.connectedPeers[c.RemotePeer()] = struct{}{}
	}
}

func (cmgr *ConnectionManager) handleNewUser(peerID peer.ID) {
	// TODO: send them all the stored messages
}

func (cmgr *ConnectionManager) receiveNewMessages(hub *api.Hub, peerID peer.ID) {
	for _, message := range hub.Messages {
		// TODO: Get rid of the channel and just store the messages directly
		cmgr.NewMessages <- message
	}
	go cmgr.BroadcastMessages(hub.Messages)
}

func InitConnectionManager() (*ConnectionManager, error) {
	cmgr := ConnectionManager{peersMutex: sync.Mutex{}}
	_self, err := libp2p.New()
	if err != nil {
		return nil, err
	}
	cmgr.host = _self
	cmgr.connectedPeers = make(map[peer.ID]struct{})
	cmgr.NewMessages = make(chan api.Message, 10)
	cmgr.peerAddrChan = make(chan peer.AddrInfo, 10)

	// register disconnect protocol
	cmgr.host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: cmgr.peerDisconnectWrapper(),
		ConnectedF:    cmgr.peerConnectWrapper(),
	})

	// register stream handlers for protocols
	cmgr.host.SetStreamHandler(protocol.ID(identify.HandshakeProtocolID),
		func(stream network.Stream) {
			identify.IdentifyHandler(stream, false, cmgr.handleNewUser)
		})

	cmgr.host.SetStreamHandler(protocol.ID(send_message.ProtocolID),
		func(stream network.Stream) {
			send_message.SendMessageHandler(stream, cmgr.receiveNewMessages)
		})

	// initialise peer discovery via mdns
	cmgr.peerAddrChan, cmgr.mDNSservice, err = discovery.InitMDNS(cmgr.host)
	if err != nil {
		return nil, err
	}

	// bin all incoming peers, they should challenge us
	go func() {
		for {
			<-cmgr.peerAddrChan
		}
	}()

	return &cmgr, nil
}

func (cmgr *ConnectionManager) BroadcastMessages(msgs []api.Message) {
	outgoing := api.Hub{
		HubName:  cmgr.self.HubName,
		HubID:    cmgr.self.HubID,
		Messages: msgs,
	}
	for peerID := range cmgr.connectedPeers {
		go send_message.SendMessages(cmgr.host,
			context.Background(), peerID, outgoing)
	}
}
