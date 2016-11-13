package server

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/gorilla/securecookie"
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
	cookie = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	)
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

const cookieName = "presence"

func (h *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := fmt.Sprint(time.Now().UnixNano())
	encoded, err := cookie.Encode(cookieName, id)
	if err != nil {
		log.Error(err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: encoded,
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

// ID is the context key used to store the connection ID
const ID = "id"

// HandlerFunc wraps an http.HandlerFunc to attach the connection ID to the HTTP request context.
func HandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var id string
		if err = cookie.Decode(cookieName, c.Value, &id); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(w, r.WithContext(context.WithValue(r.Context(), ID, id)))
	}
}
