package server

import (
	"bytes"

	"golang.org/x/net/context"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/sync/errgroup"
	"github.com/gorilla/websocket"

	pb "github.com/dichro/intangible"
	log "github.com/golang/glog"
)

type Conn struct {
	id   string
	conn *websocket.Conn
}

func NewConn(id string, conn *websocket.Conn) *Conn {
	return &Conn{
		conn: conn,
		id:   id,
	}
}

type Connecter interface {
	Connect(context.Context, string, <-chan *pb.ClientUpdate) <-chan []byte
}

func (c *Conn) Connect(ctx context.Context, room Connecter) {
	defer c.conn.Close()
	defer log.Infof("ending connect %s", c.id)
	g, ctx := errgroup.WithContext(ctx)
	in := make(chan *pb.ClientUpdate)
	g.Go(func() error { return c.pullLoop(ctx, in) })
	out := room.Connect(ctx, c.id, in)
	g.Go(func() error { return c.pushLoop(ctx, out) })
	if err := g.Wait(); err != nil {
		log.Error(err)
	}
}

func (c *Conn) pushLoop(ctx context.Context, ch <-chan []byte) error {
	defer log.Infof("ending pushLoop %s", c.id)
	for msg := range ch {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conn) pullLoop(ctx context.Context, ch chan<- *pb.ClientUpdate) error {
	defer log.Infof("ending pullLoop %s", c.id)
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}
		var rcvd pb.ClientUpdate
		if err = jsonpb.Unmarshal(bytes.NewReader(p), &rcvd); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case ch <- &rcvd:
		}
	}
}
