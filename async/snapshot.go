package async

import (
	"sync"

	"golang.org/x/net/context"
)

// Update is a composable update for some abstract object.
type Update interface {
	// Apply applies this Update to a Snapshot. Snapshot will be locked during this call (ie, Snapshot itself is modifiable), though values within Snapshot may not be, and should therefore be replaced, rather than modified.
	Apply(Snapshot) error
	// Exists returns true when the Snapshot already contains this Update.
	Exists(Snapshot) bool
}

// Snapshot is a string-keyed map of Updates.
type Snapshot map[string]Update

// SyncSnapshot allows clients to monitor changes to a Snapshot as Updates are applied.
type SyncSnapshot struct {
	// mu also protects all accesses to ring. Any call to ring that triggers a lock on its internal RWMutex is only made when holding the corresponding lock on mu.
	mu sync.RWMutex
	// Ring provides pubsub
	ring   *RingBuffer
	latest Snapshot
}

// NewSyncSnapshot returns a new SyncSnapshot configured to allow subscribers to lag by at least maxBacklog updates without requiring a Resync.
func NewSyncSnapshot(maxBacklog int) *SyncSnapshot {
	return &SyncSnapshot{
		ring:   NewRingBuffer(maxBacklog),
		latest: make(Snapshot),
	}
}

// Apply applies this Update to this SyncSnapshot.
func (s *SyncSnapshot) Apply(update Update) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := update.Apply(s.latest); err != nil {
		return err
	}
	s.ring.Push(update)
	return nil
}

// Resync returns any differences this SyncSnapshot and the client's lastKnown, and a channel that will stream future updates to this SyncSnapshot.
func (s *SyncSnapshot) Resync(ctx context.Context, lastKnown Snapshot) (updates chan []interface{}, updated Snapshot, removed []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// no updates will be delivered to Ring until we unlock s.mu, hence this channel should stay idle until Resync returns, except for the initial sentinel write from Ring.Watch.
	updates = make(chan []interface{})
	go s.ring.Watch(ctx, updates)

	updated = make(Snapshot)
	for id, latest := range s.latest {
		if !latest.Exists(lastKnown) {
			updated[id] = latest
		}
	}

	for id := range lastKnown {
		if _, ok := s.latest[id]; !ok {
			removed = append(removed, id)
		}
	}

	// wait for sentinel from ringbuffer to make sure that the channel is synchronized to current state
	<-updates
	return
}
