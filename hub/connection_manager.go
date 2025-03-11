package connection

import (
	"context"
	"log"
	"sync"

	hub_storage "github.com/khorcarol/AgentOfThings/hub/storage"
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
	log.Printf("Sending all old messages to new peer %+v\n", peerID)
	// send them all the stored messages
	msgs, _ := hub_storage.ReadMessages()
	outgoing := api.Hub{
		HubName:  cmgr.self.HubName,
		HubID:    cmgr.self.HubID,
		Messages: msgs,
	}
	go send_message.SendMessages(cmgr.host,
		context.Background(), peerID, outgoing)
}

func (cmgr *ConnectionManager) receiveNewMessages(hub *api.Hub, peerID peer.ID) {
	log.Printf("Receiving new message from peer %+v\n", peerID)
	for _, message := range hub.Messages {
		hub_storage.StoreMessage(message)
	}
	go cmgr.BroadcastMessages(hub.Messages)
}

func InitConnectionManager(hub api.Hub) (*ConnectionManager, error) {
	_self, err := libp2p.New()
	if err != nil {
		return nil, err
	}
	cmgr := ConnectionManager{
		peersMutex:     sync.Mutex{},
		host:           _self,
		connectedPeers: make(map[peer.ID]struct{}),
		NewMessages:    make(chan api.Message, 10),
		peerAddrChan:   make(chan peer.AddrInfo, 10),
		self:           hub,
	}

	// register disconnect protocol
	cmgr.host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: cmgr.peerDisconnectWrapper(),
		ConnectedF:    cmgr.peerConnectWrapper(),
	})

	// register stream handlers for protocols
	cmgr.host.SetStreamHandler(protocol.ID(identify.HandshakeProtocolID),
		func(stream network.Stream) {
			log.Printf("Responding to identify from peer %+v\n", stream.Conn().RemotePeer())
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
