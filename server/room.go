package server

import (
	"fmt"
	"log"
	"sync"
	"time"

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
func (r *Room) Place(object *pb.Object) string {
	id := r.nextID()
	object.Id = id
	buf, err := proto.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}
	r.state.Update(id, 1, buf)
	return id
}

func (r *Room) nextID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := fmt.Sprintf("%d|%d|%s", r.id, r.generation, r.name)
	r.id++
	return id
}
