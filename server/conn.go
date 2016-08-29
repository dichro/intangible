package server

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/dichro/intangible/async"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"

	pb "github.com/dichro/intangible"
	log "github.com/golang/glog"
)

var (
	sceneUpdates = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_scene_rcvd",
			Help: "number of updates received from scene",
		},
		[]string{"id"},
	)
	syncedUpdates = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_client_sent",
			Help: "number of updates synced to client",
		},
		[]string{"id"},
	)
)

func init() { prometheus.MustRegister(sceneUpdates, syncedUpdates) }

type Conn struct {
	id    string
	conn  *websocket.Conn
	close context.CancelFunc
}

func NewConn(id string, conn *websocket.Conn) *Conn {
	return &Conn{
		conn: conn,
		id:   fmt.Sprint(time.Now().UnixNano()),
	}
}

func (c *Conn) ConnectRoom(ctx context.Context, room *async.LatestSnapshot) error {
	ctx, c.close = context.WithCancel(ctx)
	defer c.close()
	defer room.Delete(c.id)
	go c.pushLoop(ctx, room)
	go c.pullLoop(ctx, room)
	<-ctx.Done()
	log.Infof("client %s disconnected", c.id)
	time.Sleep(time.Minute)
	sceneUpdates.DeleteLabelValues(c.id)
	syncedUpdates.DeleteLabelValues(c.id)
	log.Infof("client %s removed", c.id)
	return nil
}

var closed = make(chan []interface{})

func init() { close(closed) }

func (c *Conn) pushLoop(ctx context.Context, scene *async.LatestSnapshot) {
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
				log.Infof("resynchronizing %s with %d synced to client; %d pending", c.id, len(synced), len(pending))
				updateCh = scene.Resync(ctx, synced)
				updates = <-updateCh
				pending = make(async.UpdateMap, len(updates))
				// ...and fall through to process the first set of updates that result from the resync.
			}
			sceneUpdates.WithLabelValues(c.id).Add(float64(len(updates)))
			for _, u := range updates {
				switch u, ok := u.(async.Update); {
				case !ok:
					log.Error("unknown update type")
				case u.ID == c.id:
					if u.Deleted {
						log.Error("scene deleted self??")
						return
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
				if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
					log.Error(err)
					c.close()
					return
				}
				syncedUpdates.WithLabelValues(c.id).Add(float64(len(updates)))
				if update.Deleted {
					delete(synced, id)
				} else {
					synced[id] = update
				}
			}
			// TOOD(dichro): when is it better to loop over keys and delete(pending, key) rather than recreating?
			pending = make(async.UpdateMap)
		case <-ctx.Done():
			return
		}
	}
}

func (c *Conn) pullLoop(ctx context.Context, scene *async.LatestSnapshot) {
	defer c.close()
	nextVersion := uint64(1)
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}
		var rcvd pb.Object
		if err = proto.Unmarshal(p, &rcvd); err != nil {
			log.Error(err)
			continue
		}
		// TODO(dichro): validate updates
		// force ID to match assigned
		rcvd.Id = c.id
		// client can't force a bounding box
		rcvd.BoundingBox = &pb.Vector{0.3, 0.3, 0.3}
		// TODO(dichro): merge into current state! Don't rely on complete picture in every update.

		// Pre-prepare the protos for sending to any clients
		buf, err := proto.Marshal(&rcvd)
		if err != nil {
			log.Error(err)
			return
		}
		scene.Update(c.id, nextVersion, buf)
		nextVersion++
	}
}
