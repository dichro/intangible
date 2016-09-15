package server

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"

	log "github.com/golang/glog"
)

var (
	upgrader           = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	clientSubprotocols = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_client_subprotocols",
			Help: "subprotocols supported by clients",
		},
		[]string{"subprotocol"},
	)
)

func init() { prometheus.MustRegister(clientSubprotocols) }

const subprotocol = "intangible-v0"

func (rm *Room) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sps := websocket.Subprotocols(r)
	resp := make(http.Header)
	if len(sps) > 0 {
		for _, sp := range sps {
			// TODO(dichro): can labels contain bad data? Should I URI-escape?
			clientSubprotocols.WithLabelValues(sp).Inc()
			if sp == subprotocol {
				resp.Set("Sec-Websocket-Protocol", sp)
			}
		}
	} else {
		clientSubprotocols.WithLabelValues("none").Inc()
	}
	ws, err := upgrader.Upgrade(w, r, resp)
	if err != nil {
		log.Error(err)
		return
	}
	id := fmt.Sprint(time.Now().UnixNano())
	conn := NewConn(id, ws)
	conn.ConnectRoom(context.Background(), rm.state)
}
