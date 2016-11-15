package server

import (
	"errors"

	"github.com/dichro/intangible/async"

	pb "github.com/dichro/intangible"
)

var (
	// Apply errors
	ErrDifferentID = errors.New("different ID in merge")
	ErrNotObject   = errors.New("existing value in map is not a *Object")
	ErrRollback    = errors.New("can't apply older version to newer one")
)

// Update implements async.Update for async.UpdateMaps of *Updates.
type Update struct {
	version uint64
	o       *pb.Object
}

func NewUpdate(version uint64, update *pb.Object) *Update {
	return &Update{version, update}
}

func NewDelete(id string) *Update {
	return &Update{
		o: &pb.Object{
			Id:      id,
			Removed: true,
		},
	}
}

func (o *Update) Proto() *pb.Object {
	return o.o
}

// Update applies this Object to an UpdateMap that contains *Object values. If the map already contains a matching ID, a copy of the existing Object is made, updated, and re-inserted into the map.
func (o *Update) Apply(m async.Snapshot) error {
	ref, ok := m[o.o.Id]
	if !ok {
		m[o.o.Id] = o
		return nil
	}
	existing, ok := ref.(*Update)
	if !ok {
		return ErrNotObject
	}
	if o.o.Id != existing.o.Id {
		return ErrDifferentID
	}
	// removals are not versioned
	if o.o.Removed {
		delete(m, o.o.Id)
		return nil
	}
	if o.version < existing.version {
		return ErrRollback
	}
	// crude shallow copy
	obj := *existing.o
	if o.o.Position != nil {
		obj.Position = o.o.Position
	}
	if o.o.Rotation != nil {
		obj.Rotation = o.o.Rotation
	}
	if o.o.Rendering != nil {
		obj.Rendering = o.o.Rendering
	}
	x := *existing
	x.o = &obj
	x.version = o.version
	m[o.o.Id] = &x
	return nil
}

func (o *Update) Exists(m async.Snapshot) bool {
	ref, ok := m[o.o.Id]
	if !ok {
		return false
	}
	existing, ok := ref.(*Update)
	if !ok {
		return false
	}
	return o.version < existing.version
}
