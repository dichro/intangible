package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/dichro/intangible/async"
	"github.com/golang/protobuf/proto"

	pb "github.com/dichro/intangible"
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
