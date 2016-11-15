package server

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/dichro/intangible/async"
	"github.com/golang/protobuf/proto"
	"github.com/golang/sync/errgroup"

	pb "github.com/dichro/intangible"
	log "github.com/golang/glog"
)

// Room is a standalone, network-connected virtual space.
type Room struct {
	name       string
	generation int64
	state      *async.SyncSnapshot

	mu sync.Mutex
	id int
}

// NewRoom returns a new room.
func NewRoom(name string) *Room {
	return &Room{
		name:       name,
		generation: time.Now().UnixNano(),
		// TODO(dichro): how should the queue size be decided?
		state: async.NewSyncSnapshot(20),
	}
}

// Place adds the object to this Room and returns its ID.
func (r *Room) Place(object *pb.Object) *Presence {
	p := &Presence{
		ID:      r.nextID(),
		Room:    r,
		version: 1,
	}
	p.Update(object)
	return p
}

func (r *Room) nextID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := fmt.Sprintf("%d|%d|%s", r.id, r.generation, r.name)
	r.id++
	return id
}

func (r *Room) Connect(ctx context.Context, id string, updates <-chan *pb.Object) <-chan []byte {
	log.Infof("client %s connected", id)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return r.pull(id, updates) })
	ch := make(chan []byte)
	g.Go(func() error { return r.push(ctx, id, ch) })
	go func() {
		defer r.state.Apply(NewDelete(id))
		if err := g.Wait(); err != nil {
			log.Error(err)
		}
		log.Infof("client %s disconnected", id)
		time.Sleep(time.Minute)
		log.Infof("client %s removed", id)
	}()
	return ch
}

func (r *Room) pull(id string, updates <-chan *pb.Object) error {
	nextVersion := uint64(1)
	for rcvd := range updates {
		// TODO(dichro): validate updates
		// force ID to match assigned
		rcvd.Id = id
		// client can't force a bounding box
		rcvd.BoundingBox = &pb.Vector{0.3, 0.3, 0.3}
		r.state.Apply(NewUpdate(nextVersion, rcvd))
		nextVersion++
	}
	return nil
}

var closed = make(chan []interface{})

func init() { close(closed) }

func (r *Room) push(ctx context.Context, id string, ch chan<- []byte) error {
	var (
		// updateCh is automatically re-opened and re-synchronized whenever it's closed.
		updateCh = closed
		// updates received from updateCh, but not synced to client
		pending async.Snapshot
		// updates synced to client
		synced = make(async.Snapshot)
	)
	for {
		select {
		case updates, ok := <-updateCh:
			if !ok {
				// Throw away our pending list and re-build it by synchronizing against what we've already sent to the client.
				log.Infof("resynchronizing %s with %d synced to client; %d pending", id, len(synced), len(pending))
				ch, add, remove := r.state.Resync(ctx, synced)
				updateCh = ch
				for _, id := range remove {
					updates = append(updates, NewDelete(id))
				}
				for _, u := range add {
					updates = append(updates, u)
				}
				pending = make(async.Snapshot, len(updates))
				// ...and fall through to process the updates that result from the resync.
			}
			for _, u := range updates {
				if u, ok := u.(async.Update); ok {
					u.Apply(pending)
				} else {
					log.Error("unknown update type")
				}
			}
			for _, u := range pending {
				// TODO(dichro): limit elapsed time/number of loop iterations to reduce backlog on updateCh?
				update, ok := u.(*Update)
				if !ok {
					log.Error("not a *server.Update payload in update.Data")
					continue
				}
				buf, err := proto.Marshal(update.Proto())
				if err != nil {
					log.Error(err)
					continue
				}
				ch <- buf
				u.Apply(synced)
			}
			// TOOD(dichro): when is it better to loop over keys and delete(pending, key) rather than recreating?
			pending = make(async.Snapshot)
		case <-ctx.Done():
			return nil
		}
	}
}

// Presence handles the presence of an object in a room.
type Presence struct {
	ID     string
	Object *pb.Object
	Room   *Room

	mu      sync.Mutex
	removed bool
	version uint64
}

func (p *Presence) nextVersion() uint64 {
	version := p.version
	p.version++
	return version
}

// Update updates the state of an object.
func (p *Presence) Update(update *pb.Object) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.removed {
		// TODO(dichro): error? log? panic?
		return
	}
	p.Object = update
	update.Id = p.ID
	p.Room.state.Apply(NewUpdate(p.nextVersion(), update))
}

// Remove removes this object from its room.
func (p *Presence) Remove() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Room.state.Apply(NewDelete(p.ID))
	p.removed = true
}
