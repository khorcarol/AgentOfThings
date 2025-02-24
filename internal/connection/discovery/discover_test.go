package discovery

import (
	"fmt"
	"testing"

	libp2p "github.com/libp2p/go-libp2p"
)

func TestCanFind(t *testing.T) {
	h1, err := libp2p.New()
	if err != nil {
		t.Errorf("%s", err)
	}
	c1, s1, err := InitMDNS(h1)
	if err != nil {
		t.Errorf("%s", err)
	}

	h2, err := libp2p.New()
	if err != nil {
		t.Errorf("%s", err)
	}
	c2, s2, err := InitMDNS(h2)
	if err != nil {
		t.Errorf("%s", err)
	}

	p1 := <-c1
	p2 := <-c2
	CloseMDNS(s1)
	CloseMDNS(s2)

	fmt.Printf("Peer 1 ID: %v\nPeer 2 ID: %v\n", p1.ID, p2.ID)
	if p1.ID != h2.ID() && p2.ID != h1.ID() {
		t.Fail()
	}
}
