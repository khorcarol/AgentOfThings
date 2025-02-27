package connection

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection/discovery"
	"github.com/khorcarol/AgentOfThings/internal/connection/protocol/handshake/peer_to_user"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type peerLevel int

const (
	Peer = iota
	User
	Friend
)

type ConnectionManager struct {
	mDNSservice    mdns.Service
	peerAddrChan   <-chan peer.AddrInfo
	host           host.Host
	peersMutex     sync.Mutex
	connectedPeers map[peer.ID]peerLevel
	uuids          map[uuid.UUID]peer.ID

	Self api.User
	// B->M, sends a new discovered user
	IncomingUsers chan api.User
	// B->M, sends the response to the fried request (potentially with data)
	IncomingFriendResponse chan api.FriendResponse
	// B->M, sends an external friend request to respond to
	IncomingFriendRequest chan api.Friend
}

func (cmgr *ConnectionManager) peerDisconnectWrapper() func(network.Network, network.Conn) {
	return func(n network.Network, c network.Conn) {
		defer cmgr.peersMutex.Unlock()
		cmgr.peersMutex.Lock()
		delete(cmgr.connectedPeers, c.RemotePeer())
	}
}

func (cmgr *ConnectionManager) peerConnectWrapper() func(network.Network, network.Conn) {
	return func(n network.Network, c network.Conn) {
		cmgr.peersMutex.Lock()
		defer cmgr.peersMutex.Unlock()
		cmgr.connectedPeers[c.RemotePeer()] = Peer
	}
}

func initConnectionManager() (*ConnectionManager, error) {
	cmgr := ConnectionManager{}
	_self, err := libp2p.New()
	if err != nil {
		return nil, err
	}
	cmgr.host = _self
	cmgr.connectedPeers = make(map[peer.ID]peerLevel)
	cmgr.IncomingUsers = make(chan api.User)
	cmgr.IncomingFriendResponse = make(chan api.FriendResponse)
	cmgr.IncomingFriendRequest = make(chan api.Friend)

	// register disconnect protocol
	cmgr.host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: cmgr.peerDisconnectWrapper(),
		ConnectedF:    cmgr.peerConnectWrapper(),
	})

	// TODO: register stream handlers for protocols
	cmgr.host.SetStreamHandler(protocol.ID(peer_to_user.HandshakeProtocolID),
		func(stream network.Stream) {
			peer_to_user.HandshakeHandler(stream, cmgr.addIncomingUser, &cmgr.Self)
		})

	// initialise peer discovery via mdns
	cmgr.peerAddrChan, cmgr.mDNSservice, err = discovery.InitMDNS(cmgr.host)
	if err != nil {
		return nil, err
	}
	return &cmgr, nil
}

func (cmgr *ConnectionManager) addIncomingUser(msg *api.User, id peer.ID) {
	if msg == nil {
		return
	}
	cmgr.peersMutex.Lock()
	defer cmgr.peersMutex.Unlock()
	cmgr.uuids[msg.UserID.Address] = id
	cmgr.connectedPeers[id] = User
	cmgr.IncomingUsers <- *msg
}

func (cmgr *ConnectionManager) SendFriendRequest(user api.User, data api.Friend) error {
	// start a new stream with friend request protocol
	return nil
}

func (cmgr *ConnectionManager) SendFriendResponse(user api.User, data api.FriendResponse) error {
	return nil
}

func (cmgr *ConnectionManager) waitOnPeer(wg *sync.WaitGroup) {
	wg.Add(1)
	peerAddr := <-cmgr.peerAddrChan
	go cmgr.connectToPeer(peerAddr, wg)
}

func (cmgr *ConnectionManager) peerToUserHandshake(peerAddr peer.AddrInfo) error {
	// check whether we are actually still connected to the peer
	if _, ok := cmgr.connectedPeers[peerAddr.ID]; !ok {
		return nil
	}
	// carry out handshake
	msg, err := peer_to_user.InitiateHandshake(cmgr.host, context.Background(), peerAddr.ID, cmgr.Self)
	if err != nil {
		return err
	}
	if msg == nil {
		return nil
	}

	// Add user to set of currently connected users and send to middle
	cmgr.peersMutex.Lock()
	defer cmgr.peersMutex.Unlock()
	cmgr.uuids[msg.UserID.Address] = peerAddr.ID
	cmgr.connectedPeers[peerAddr.ID] = User
	cmgr.IncomingUsers <- *msg
	return nil
}

func (cmgr *ConnectionManager) connectToPeer(peerAddr peer.AddrInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	cmgr.peersMutex.Lock()
	defer cmgr.peersMutex.Unlock()
	// check if we are already connected to this peer
	if _, ok := cmgr.connectedPeers[peerAddr.ID]; ok {
		return
	}
	// open connection to peer
	if err := cmgr.host.Connect(context.Background(), peerAddr); err != nil {
		log.Fatal("Failed to connect to new peer", err)
	}
	cmgr.connectedPeers[peerAddr.ID] = Peer

	// handshake to promote peer to user
	go cmgr.peerToUserHandshake(peerAddr)
}
