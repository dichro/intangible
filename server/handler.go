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
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
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

// WebsocketHandler handles websocket connections from clients.
type WebsocketHandler struct {
	c Connecter
}

func NewWebsocketHandler(room Connecter) *WebsocketHandler {
	return &WebsocketHandler{
		c: room,
	}
}

func (h *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := fmt.Sprint(time.Now().UnixNano())
	http.SetCookie(w, &http.Cookie{
		Name:  "presence",
		Value: id,
		Path:  "/",
	})
	resp := w.Header()
	sps := websocket.Subprotocols(r)
	if len(sps) > 0 {
		for _, sp := range sps {
			if sp == subprotocol {
				clientSubprotocols.WithLabelValues(sp).Inc()
				resp.Set("Sec-Websocket-Protocol", sp)
			}
		}
	} else {
		clientSubprotocols.WithLabelValues("none").Inc()
	}
	ws, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Error(err)
		return
	}
	conn := NewConn(id, ws)
	conn.Connect(context.Background(), h.c)
}
