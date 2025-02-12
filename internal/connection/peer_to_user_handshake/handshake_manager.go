package peer_to_user_handshake

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	handshakeValidityDuration = 1 * time.Hour
	handshakeFailureTimeout   = 10 * time.Minute // Prevents rediscovering the same "new" peer over and over and flooding them.
	maxHandshakeAttempts      = 3
	retryDelay                = 2 * time.Second
)

type handshakeRecord struct {
	lastAttempt time.Time
	success     bool
	failedCount int
}

// HandshakeRegistry is a table of peer.ID -> [handshakeRecord] mappings to track when handshakes were last attempted.
type HandshakeRegistry struct {
	mutex   sync.Mutex
	records map[peer.ID]*handshakeRecord
}

func NewHandshakeRegistry() *HandshakeRegistry {
	return &HandshakeRegistry{
		records: make(map[peer.ID]*handshakeRecord),
	}
}

// ShouldHandshake returns true if a handshake for the given peer should be attempted.
// It is based on the values in the [HandshakeRegistry].
func (registry *HandshakeRegistry) ShouldHandshake(peerID peer.ID) bool {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	rec, exists := registry.records[peerID]
	if !exists {
		return true
	}
	now := time.Now()
	if rec.success {
		// To not keep trying with successful peers.
		return now.Sub(rec.lastAttempt) > handshakeValidityDuration
	}
	// To not keep trying with failed peers.
	return now.Sub(rec.lastAttempt) > handshakeFailureTimeout
}

func (registry *HandshakeRegistry) RecordSuccess(peerID peer.ID) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	registry.records[peerID] = &handshakeRecord{
		lastAttempt: time.Now(),
		success:     true,
		failedCount: 0,
	}
}

func (registry *HandshakeRegistry) RecordFailure(peerID peer.ID) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	if rec, exists := registry.records[peerID]; exists {
		rec.lastAttempt = time.Now()
		rec.success = false
		rec.failedCount++
	} else {
		registry.records[peerID] = &handshakeRecord{
			lastAttempt: time.Now(),
			success:     false,
			failedCount: 1,
		}
	}
}

// HandshakeManager coordinates handshake attempts on new peers.
type HandshakeManager struct {
	Registry         *HandshakeRegistry
	HandshakeService *HandshakeService
	LocalPeerID      peer.ID
}

func NewHandshakeManager(registry *HandshakeRegistry, handshakeService *HandshakeService, localPeerID peer.ID) *HandshakeManager {
	return &HandshakeManager{
		Registry:         registry,
		HandshakeService: handshakeService,
		LocalPeerID:      localPeerID,
	}
}

// MaybeHandshake attempts to perform a handshake with a newly discovered peer.
// It first checks the registry to ensure that either no recent attempt exists or that
// enough time has passed since the last attempt (whether that attempt was a success or failure),
// and only if ShouldInitiate indicates that our local side should initiate.
// It then retries on failure up to maxHandshakeAttempts. If all attempts fail,
// it records the failure so that we don’t try again until handshakeFailureTimeout expires.
func (manager *HandshakeManager) MaybeHandshake(ctx context.Context, remotePeerID peer.ID) {
	if !manager.Registry.ShouldHandshake(remotePeerID) {
		log.Printf("handshake_manager: peer %s was attempted too recently, skipping handshake", remotePeerID)
		return
	}

	if !ShouldInitiate(manager.LocalPeerID, remotePeerID) {
		log.Printf("handshake_manager: not initiating handshake with peer %s", remotePeerID)
		return
	}

	var lastErr error
	for attempt := 1; attempt <= maxHandshakeAttempts; attempt++ {
		log.Printf("handshake_manager: handshake attempt %d with peer %s", attempt, remotePeerID)
		handshakeCtx, cancel := context.WithTimeout(ctx, HandshakeTimeout())
		// Each attempt gets its own cancel context.
		func() {
			defer cancel()
			receivedMsg, err := manager.HandshakeService.InitiateHandshake(handshakeCtx, remotePeerID)
			if err == nil {
				log.Printf("handshake_manager: handshake with peer %s succeeded: %+v", remotePeerID, receivedMsg)
				manager.Registry.RecordSuccess(remotePeerID)
				lastErr = nil
				return
			}
			lastErr = err
			log.Printf("handshake_manager: handshake with peer %s failed on attempt %d: %v", remotePeerID, attempt, err)
		}()
		if lastErr == nil {
			return
		}
		time.Sleep(retryDelay)
	}
	log.Printf("handshake_manager: handshake with peer %s failed after %d attempts: %v", remotePeerID, maxHandshakeAttempts, lastErr)
	manager.Registry.RecordFailure(remotePeerID)
}
