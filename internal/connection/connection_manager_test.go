package connection

import (
	"fmt"
	"sync"
	"testing"
)

func TestCanConnect(t *testing.T) {
	cmgr1, err := InitConnectionManager()
	if err != nil {
		t.Errorf("Failed with error %s", err)
	}
	cmgr2, err := InitConnectionManager()
	if err != nil {
		t.Errorf("Failed with error %s", err)
	}
	var wg sync.WaitGroup
	cmgr1.waitOnPeer(&wg)
	wg.Wait()
	if len(cmgr1.connectedPeers) == 0 || len(cmgr1.connectedPeers) != len(cmgr2.connectedPeers) {
		t.Fail()
	}
	fmt.Printf("Host1 len %+v, Host2 len %+v\n", len(cmgr1.connectedPeers), len(cmgr2.connectedPeers))
}
