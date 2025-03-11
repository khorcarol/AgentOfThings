package connection

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection/discovery"
	"github.com/khorcarol/AgentOfThings/internal/connection/protocol/handshake/peer_to_user"
	"github.com/khorcarol/AgentOfThings/internal/connection/protocol/handshake/user_to_friend"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/vishalkuo/bimap"
)

type peerLevel int

type ConnectionManager struct {
	mDNSservice    mdns.Service
	peerAddrChan   <-chan peer.AddrInfo
	host           host.Host
	peersMutex     sync.Mutex
	connectedPeers map[peer.ID]struct{}
	uuids          bimap.BiMap[uuid.UUID, peer.ID]

	// B->M, sends a new discovered user
	IncomingUsers chan api.User
	// B->M, sends an external friend request to respond to
	IncomingFriendRequest chan api.FriendRequest
	// B->M, informs middle that a peer has disconnected
	PeerDisconnections chan uuid.UUID
	// B->M, informs middle of new messages
	NewMessages chan api.Hub
}

func (cmgr *ConnectionManager) peerDisconnectWrapper() func(network.Network, network.Conn) {
	return func(n network.Network, c network.Conn) {
		cmgr.peersMutex.Lock()
		defer cmgr.peersMutex.Unlock()
		log.Printf("Peer %+v disconnected\n", c.RemotePeer())
		uuid, ok := cmgr.uuids.GetInverse(c.RemotePeer())
		if ok {
			cmgr.PeerDisconnections <- uuid
		}
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

func initConnectionManager() (*ConnectionManager, error) {
	cmgr := ConnectionManager{peersMutex: sync.Mutex{}}
	_self, err := libp2p.New()
	if err != nil {
		return nil, err
	}
	cmgr.host = _self
	cmgr.connectedPeers = make(map[peer.ID]struct{})
	cmgr.uuids = *bimap.NewBiMap[uuid.UUID, peer.ID]()
	cmgr.IncomingUsers = make(chan api.User, 10)
	cmgr.IncomingFriendRequest = make(chan api.FriendRequest, 10)
	cmgr.PeerDisconnections = make(chan uuid.UUID, 10)

	// register disconnect protocol
	cmgr.host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: cmgr.peerDisconnectWrapper(),
		ConnectedF:    cmgr.peerConnectWrapper(),
	})

	// register stream handlers for protocols
	cmgr.host.SetStreamHandler(protocol.ID(peer_to_user.HandshakeProtocolID),
		func(stream network.Stream) {
			peer_to_user.HandshakeHandler(stream, cmgr.addIncomingUser)
		})
	cmgr.host.SetStreamHandler(protocol.ID(user_to_friend.FriendRequestProtocolID),
		func(stream network.Stream) {
			user_to_friend.FriendRequestHandler(stream, func(f *api.FriendRequest, pid peer.ID) {
				// Pass friend data to middle.
				cmgr.IncomingFriendRequest <- *f
			})
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
	cmgr.uuids.Insert(msg.UserID.Address, id)
	cmgr.connectedPeers[id] = struct{}{}
	cmgr.IncomingUsers <- *msg
}

// SendFriendRequest sends a friend request by calling the friend protocol layer.
// Application logic (i.e. middle) should handle if this friend is to be displayed or stored.
func (cmgr *ConnectionManager) SendFriendRequest(user api.User, data api.FriendRequest) error {
	peerID, ok := cmgr.uuids.Get(user.UserID.Address)
	if !ok {
		return nil
	}

	return user_to_friend.SendFriendData(cmgr.host, context.Background(), peerID, data)
}

func (cmgr *ConnectionManager) SendMessage(HubID api.ID, message api.Message) error {
	return nil
}

func (cmgr *ConnectionManager) waitOnPeer(wg *sync.WaitGroup) {
	wg.Add(1)
	peerAddr := <-cmgr.peerAddrChan
	go cmgr.connectToPeer(peerAddr, wg)
}

func (cmgr *ConnectionManager) peerToUserHandshake(peerAddr peer.AddrInfo, wg *sync.WaitGroup) error {
	defer wg.Done()
	// check whether we are actually still connected to the peer
	if _, ok := cmgr.connectedPeers[peerAddr.ID]; !ok {
		return nil
	}
	// carry out handshake
	msg, err := peer_to_user.InitiateHandshake(cmgr.host, context.Background(), peerAddr.ID)
	if err != nil {
		return err
	}
	if msg == nil {
		return nil
	}

	// Add user to set of currently connected users and send to middle
	cmgr.peersMutex.Lock()
	defer cmgr.peersMutex.Unlock()
	cmgr.uuids.Insert(msg.UserID.Address, peerAddr.ID)
	cmgr.connectedPeers[peerAddr.ID] = struct{}{}
	cmgr.IncomingUsers <- *msg
	return nil
}

func (cmgr *ConnectionManager) connectToPeer(peerAddr peer.AddrInfo, wg *sync.WaitGroup) {
	cmgr.peersMutex.Lock()
	// check if we are already connected to this peer
	if _, ok := cmgr.connectedPeers[peerAddr.ID]; ok {
		wg.Done()
		cmgr.peersMutex.Unlock()
		return
	}
	cmgr.peersMutex.Unlock()
	// open connection to peer
	if err := cmgr.host.Connect(context.Background(), peerAddr); err != nil {
		log.Print("Failed to connect to new peer", err)
	}

	// handshake to promote peer to user
	go cmgr.peerToUserHandshake(peerAddr, wg)
}

func (cmgr *ConnectionManager) StartDiscovery() {
	go func() {
		var wg sync.WaitGroup
		for peerAddr := range cmgr.peerAddrChan {
			if shouldHandshake(peerAddr.ID, cmgr.host.ID()) {
				wg.Add(1)
				go cmgr.connectToPeer(peerAddr, &wg)
			}
		}
	}()
}

func shouldHandshake(peerID peer.ID, ourID peer.ID) bool {
	return peerID < ourID
}
