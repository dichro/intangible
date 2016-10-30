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
	state      *async.LatestSnapshot

	mu sync.Mutex
	id int
}

// NewRoom returns a new room.
func NewRoom(name string) *Room {
	return &Room{
		name:       name,
		generation: time.Now().UnixNano(),
		// TODO(dichro): how should the queue size be decided?
		state: async.NewLatestSnapshot(20),
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

// Watch returns a channel that monitors changes in the Room. The channel may be closed unexpectedly.
func (r *Room) Watch(ctx context.Context) <-chan []interface{} {
	return r.state.Resync(ctx, nil)
}

func (r *Room) Connect(ctx context.Context, id string, updates <-chan *pb.Object) <-chan []byte {
	log.Infof("client %s connected", id)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return r.pull(id, updates) })
	ch := make(chan []byte)
	g.Go(func() error { return r.push(ctx, id, ch) })
	go func() {
		defer r.state.Delete(id)
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
		// TODO(dichro): merge into current state! Don't rely on complete picture in every update.

		// Pre-prepare the protos for sending to any clients
		buf, err := proto.Marshal(rcvd)
		if err != nil {
			return err
		}
		r.state.Update(id, nextVersion, buf)
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
		pending async.UpdateMap
		// updates synced to client
		synced = make(async.UpdateMap)
	)
	for {
		select {
		case updates, ok := <-updateCh:
			if !ok {
				// Throw away our pending list and re-build it by synchronizing against what we've already sent to the client.
				log.Infof("resynchronizing %s with %d synced to client; %d pending", id, len(synced), len(pending))
				updateCh = r.state.Resync(ctx, synced)
				updates = <-updateCh
				pending = make(async.UpdateMap, len(updates))
				// ...and fall through to process the first set of updates that result from the resync.
			}
			for _, u := range updates {
				switch u, ok := u.(async.Update); {
				case !ok:
					log.Error("unknown update type")
				case id == u.ID:
					if u.Deleted {
						return fmt.Errorf("scene deleted self (%q)??", id)
					}
					// don't bother returning client's own updates.
					synced[u.ID] = u
				case u.Deleted:
					// don't pass through deletes if we haven't already synced the object
					if _, ok := synced[u.ID]; !ok {
						delete(pending, u.ID)
						continue
					}
					fallthrough
				case u.Version != synced[u.ID].Version:
					pending[u.ID] = u
				}
			}
			for id, update := range pending {
				// TODO(dichro): limit elapsed time/number of loop iterations to reduce backlog on updateCh?
				var msg []byte
				if update.Deleted {
					var err error
					// TODO(dichro): pre-generate msg in update.Data where possible
					msg, err = proto.Marshal(&pb.Object{Id: id, Removed: true})
					if err != nil {
						log.Error(err)
						continue
					}
				} else {
					msg, ok = update.Data.([]byte)
					if !ok {
						log.Error("not a []byte payload in update.Data")
						continue
					}
				}
				ch <- msg
				if update.Deleted {
					delete(synced, id)
				} else {
					synced[id] = update
				}
			}
			// TOOD(dichro): when is it better to loop over keys and delete(pending, key) rather than recreating?
			pending = make(async.UpdateMap)
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
	buf, err := proto.Marshal(update)
	if err != nil {
		log.Fatal(err)
	}
	p.Room.state.Update(p.ID, p.nextVersion(), buf)
}

// Remove removes this object from its room.
func (p *Presence) Remove() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Room.state.Delete(p.ID)
	p.removed = true
}
