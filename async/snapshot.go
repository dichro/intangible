package async

import (
	"sync"

	"golang.org/x/net/context"
)

type Update struct {
	ID      string
	Version uint64
	Deleted bool
	Data    interface{}
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

// TODO(dichro): consider Set(), which requires that id doesn't exist; and change Update to require that it does.

// Update stores this Versioned in this snapshot. No version checking is performed: caller must ensure that Update is only called in strictly incrementing version order.
func (s *LatestSnapshot) Update(id string, version uint64, data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latest[id] = Update{
		ID:      id,
		Version: version,
		Data:    data,
	}
	s.ring.Push(s.latest[id])
}

func (s *LatestSnapshot) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.latest, id)
	s.ring.Push(makeDelete(id))
}

func makeDelete(id string) Update {
	return Update{
		ID:      id,
		Deleted: true,
	}
}

// Resync returns a channel that will stream all differences between client's latest and LatestSnapshot.
func (s *LatestSnapshot) Resync(ctx context.Context, lastKnown UpdateMap) chan []interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// no updates will be delivered to Ring until we unlock s.mu, hence this channel should stay idle until Resync returns, except for the initial sentinel write from Ring.Watch.
	// TODO(dichro): should this be a sync channel instead? a buffer of 1 requires the writer to swap out the contents with updates when the reader is behind to minimize latency, but that would then mess with the first write from Resync. A buffer of 1 allows our first write to stay in the channel without holding s.mu waiting on the reader.
	ch := make(chan []interface{}, 1)
	go s.ring.Watch(ctx, ch)

	// TODO(dichro): 100% of client's latest is unlikely to be valid, else why would they be here?
	updates := make([]interface{}, 0, len(s.latest)-len(lastKnown))

	for id := range lastKnown {
		if _, ok := s.latest[id]; !ok {
			updates = append(updates, makeDelete(id))
		}
	}
	for id, latest := range s.latest {
		if lastKnown[id].Version != latest.Version {
			updates = append(updates, latest)
		}
	}

	// Drain Ring.Watch's sentinel, replace it with our updates. This is a tiny bit dangerous: if ctx has been canceled, then Ring.Watch might start racing us after this read to close ch. But we need to wait for the first read from it before unlocking, lest we race with updates instead. Worst case is a panic on a goroutine that has already canceled its context.
	<-ch
	ch <- updates
	return ch
}
