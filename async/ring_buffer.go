package async

import (
	"sync"

	"golang.org/x/net/context"
)

// RingBuffer is a fixed size queue that signals new additions to it via channels.
type RingBuffer struct {
	mu         sync.RWMutex
	cond       *sync.Cond
	generation int
	ring       []interface{}
	head, size int
}

// NewRingBuffer returns a ring buffer with a fixed size. Any subscribing channel that develops a backlog in excess of this fixed size will be closed by Ring.
func NewRingBuffer(size int) *RingBuffer {
	r := &RingBuffer{
		ring: make([]interface{}, size),
		size: size,
	}
	r.cond = sync.NewCond(r.mu.RLocker())
	return r
}

// Push adds a new item to the Ring
func (r *RingBuffer) Push(item interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ring[r.head] = item
	r.generation++
	r.head++
	if r.head >= r.size {
		r.head = 0
	}
	r.cond.Broadcast()
}

// Peek returns the contents of the ring buffer from a specific sequence number onwards.
func (r *RingBuffer) peek(ctx context.Context, fromSeq int) (toSeq int, data []interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for r.generation == fromSeq {
		// no data buffered
		select {
		case <-ctx.Done():
			return fromSeq, nil
		default:
			r.cond.Wait()
		}
	}
	missed := r.generation - fromSeq
	if missed > r.size {
		// The ring buffer has wrapped since fromSeq. Fail.
		return r.generation, nil
	}
	data = make([]interface{}, 0, missed)
	tail := fromSeq % r.size
	head := tail + missed
	if head > r.size {
		head -= r.size
		data = append(data, r.ring[tail:]...)
		data = append(data, r.ring[:head]...)
	} else {
		data = append(data, r.ring[tail:head]...)
	}
	return r.generation, data
}

// Watch sends all updates to the RingBuffer down the supplied channel. The first item sent down the channel is a nil sentinel value, receipt of which signals that the channel is synchronized with future updates. If the backlog for the channel grows larger than the configured maximum backlog, the channel will be closed.
func (r *RingBuffer) Watch(ctx context.Context, ch chan<- []interface{}) {
	defer close(ch)
	r.mu.RLock()
	seq := r.generation
	r.mu.RUnlock()
	ch <- nil

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// TODO(dichro): go r.forceBroadcastIfNotClosed, to avoid this goroutine blocking eternally on an idle RingBuffer after a ctx.Close()

	for {
		var data []interface{}
		seq, data = r.peek(ctx, seq)
		if data == nil {
			return
		}
		// the order of these is unpredictable; in particular, note that ch <- data may be executed even after ctx.Close.
		select {
		case ch <- data:
		case <-ctx.Done():
			return
		}
	}
}
