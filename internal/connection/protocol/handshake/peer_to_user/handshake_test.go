package peer_to_user

import (
	"bytes"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/api/interests"
	"github.com/libp2p/go-libp2p/core/peer"
)

// DummyStream is a simple in-memory stream that implements the io.Reader/io.Writer interface.
type DummyStream struct {
	io.ReadWriter
	deadline time.Time
	closed   bool
	mu       sync.Mutex
}

func (ds *DummyStream) SetDeadline(t time.Time) error {
	ds.mu.Lock()
	ds.deadline = t
	ds.mu.Unlock()
	return nil
}

func (ds *DummyStream) Close() error {
	ds.mu.Lock()
	ds.closed = true
	ds.mu.Unlock()
	return nil
}

func (ds *DummyStream) Reset() error {
	return nil
}

func (ds *DummyStream) CloseWrite() error {
	return nil
}

// TestEncodeDecode verifies that we can encode a handshake message to a stream
// and then decode it back.
func TestEncodeDecode(t *testing.T) {
	origMsg := &api.User{
		Interests: []api.Interest{
			{
				Category:    interests.Sport,
				Description: "Test interest",
			},
		},
	}

	var buf bytes.Buffer
	if err := encodeToStream(&buf, origMsg); err != nil {
		t.Fatalf("encodeToStream failed: %v", err)
	}

	var decodedMsg api.User
	if err := decodeFromStream(&buf, &decodedMsg); err != nil {
		t.Fatalf("decodeFromStream failed: %v", err)
	}

	if len(decodedMsg.Interests) != 1 {
		t.Errorf("expected 1 interest; got %d", len(decodedMsg.Interests))
	}
	if decodedMsg.Interests[0].Category != origMsg.Interests[0].Category ||
		decodedMsg.Interests[0].Description != origMsg.Interests[0].Description {
		t.Errorf("decoded message does not match the original")
	}
}

// TestShouldInitiate exercises ShouldInitiate to ensure that the lexicographical
// comparison of peer IDs is correct.
func TestShouldInitiate(t *testing.T) {
	peerA := peer.ID("A")
	peerB := peer.ID("B")

	if !ShouldInitiate(peerA, peerB) {
		t.Errorf("expected A < B to return true")
	}
	if ShouldInitiate(peerB, peerA) {
		t.Errorf("expected B < A to return false")
	}
}
