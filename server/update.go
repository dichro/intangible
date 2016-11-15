package server

import (
	"errors"

	"github.com/dichro/intangible/async"

	pb "github.com/dichro/intangible"
)

var (
	// ErrDifferentID is returned when the key for an Update in a Snapshot doesn't match the Update's ID.
	ErrDifferentID = errors.New("different ID in merge")
	// ErrNotObject is returned when a Snapshot unexpectedly contains something other than a *Update.
	ErrNotObject = errors.New("existing value in map is not a *Update")
	// ErrRollback is returned when an Update is applied to a Snapshot that already contains an Update for the same key with a later version.
	ErrRollback = errors.New("can't apply older version to newer one")
)

// Update implements async.Update for async.UpdateMaps of *Updates.
type Update struct {
	version uint64
	obj     *pb.Object
}

// NewUpdate returns a new Update wrapping the provided update protobuf.
func NewUpdate(version uint64, update *pb.Object) *Update {
	return &Update{version, update}
}

// NewDelete returns an Update that destroys the provided ID.
func NewDelete(id string) *Update {
	return &Update{
		obj: &pb.Object{
			Id:      id,
			Removed: true,
		},
	}
}

// ID returns the object ID that this Update applies to.
func (u *Update) ID() string { return u.obj.Id }

// Proto returns this Update as a protobuf.
func (u *Update) Proto() *pb.Object {
	return u.obj
}

// Apply applies this Update to an async.Snapshot. If the map already contains a matching ID, a copy of the existing Update is made, updated, and re-inserted into the map.
func (u *Update) Apply(m async.Snapshot) error {
	id := u.ID()
	ref, ok := m[id]
	if !ok {
		m[id] = u
		return nil
	}
	existing, ok := ref.(*Update)
	if !ok {
		return ErrNotObject
	}
	if id != existing.obj.Id {
		return ErrDifferentID
	}
	// removals are not versioned
	if u.obj.Removed {
		delete(m, id)
		return nil
	}
	if u.version < existing.version {
		return ErrRollback
	}
	// crude shallow copy
	obj := *existing.obj
	if u.obj.Position != nil {
		obj.Position = u.obj.Position
	}
	if u.obj.Rotation != nil {
		obj.Rotation = u.obj.Rotation
	}
	if u.obj.Rendering != nil {
		obj.Rendering = u.obj.Rendering
	}
	x := *existing
	x.obj = &obj
	x.version = u.version
	m[id] = &x
	return nil
}

// Exists returns true if the provided Snapshot contains an Update with the same key and at least as old a version as this Update.
func (u *Update) Exists(m async.Snapshot) bool {
	ref, ok := m[u.obj.Id]
	if !ok {
		return false
	}
	existing, ok := ref.(*Update)
	if !ok {
		return false
	}
	return u.version < existing.version
}
