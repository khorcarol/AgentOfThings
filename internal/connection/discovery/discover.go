package discovery

import (
	"errors"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

const (
	// Size of the channel buffer
	chanSize = 10
)

func InitMDNS(peerHost host.Host) (<-chan peer.AddrInfo, mdns.Service, error) {
	serviceTag := "group-alpha/agent-of-things"
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo, chanSize)

	service := mdns.NewMdnsService(peerHost, serviceTag, n)
	if err := service.Start(); err != nil {
		return nil, nil, err
	}

	return (<-chan peer.AddrInfo)(n.PeerChan), service, nil
}

func CloseMDNS(service mdns.Service) error {
	if service == nil {
		return errors.New("discovery: mDNS service not initialised")
	}
	if err := service.Close(); err != nil {
		return err
	}
	return nil
}
