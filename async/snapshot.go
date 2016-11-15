package async

import (
	"sync"

	"golang.org/x/net/context"
)

type Update interface {
	Update(UpdateMap) error
	Exists(UpdateMap) bool
}

type UpdateMap map[string]Update

// LatestSnapshot maintains a snapshot of latest Versioneds.
type LatestSnapshot struct {
	// mu also protects all accesses to ring. Any call to ring that triggers a lock on its internal RWMutex is only made when holding the corresponding lock on mu.
	mu sync.RWMutex
	// Ring provides pubsub
	ring   *RingBuffer
	latest UpdateMap
}

// NewLatestSnapshot returns a new LatestSnapshot configured to allow subscribers to lag by at least maxBacklog updates without requiring a Resync.
func NewLatestSnapshot(maxBacklog int) *LatestSnapshot {
	return &LatestSnapshot{
		ring:   NewRingBuffer(maxBacklog),
		latest: make(UpdateMap),
	}
}

// Update stores this Versioned in this snapshot. No version checking is performed: caller must ensure that Update is only called in strictly incrementing version order.
func (s *LatestSnapshot) Update(id string, update Update) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := update.Update(s.latest); err != nil {
		return err
	}
	s.ring.Push(update)
	return nil
}

// Resync returns a channel that will stream all differences between client's latest and LatestSnapshot.
func (s *LatestSnapshot) Resync(ctx context.Context, lastKnown UpdateMap) (updates chan []interface{}, updated UpdateMap, removed []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// no updates will be delivered to Ring until we unlock s.mu, hence this channel should stay idle until Resync returns, except for the initial sentinel write from Ring.Watch.
	updates = make(chan []interface{})
	go s.ring.Watch(ctx, updates)

	updated = make(UpdateMap)
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

	return
}
