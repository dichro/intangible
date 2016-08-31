package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"

	pb "github.com/dichro/intangible"
	log "github.com/golang/glog"
)

var (
	server         = flag.String("server", "ws://localhost:8080/ws", "Server to connect to")
	updateInterval = flag.Duration("update_interval", 100*time.Millisecond, "Update interval to server")
	period         = flag.Duration("period", 30*time.Second, "Period of orbit")
	radius         = flag.Float64("radius", 1, "Radius of orbit")
	port           = flag.Int("port", 7777, "listening port for API requests")
	apiAddress     = flag.String("api_address", fmt.Sprintf("localhost:%d", *port), "published address for api clients")
)

func main() {
	flag.Parse()
	conn, _, err := (&websocket.Dialer{}).Dial(*server, nil)
	if err != nil {
		log.Exit(err)
	}
	var (
		tick   = time.Tick(*updateInterval)
		stop   = make(chan struct{})
		start  = make(chan struct{})
		update = &pb.Object{
			Id: "self",
			Position: &pb.Vector{
				Y: 1,
			},
			Api: []*pb.API{
				api(stop, "stop", "stop orbiter"),
				api(start, "start", "start orbiter"),
			},
			// Moon skin from http://tf3dm.com/3d-model/moon-17150.html
			// by Nicola Cornolti
			Rendering: &pb.Rendering{
				Mesh: &pb.Mesh{
					SourceUri: "http://intangible-gallery.s3-website-us-west-1.amazonaws.com/moon/moon.obj",
				},
				Texture: &pb.Texture{
					SourceUri: "https://s3-us-west-1.amazonaws.com/intangible-gallery/moon/MoonMap2_2500x1250.jpg",
				},
			},
		}
		pos     = 0.0
		running = true
	)
	go http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	for {
		select {
		case <-stop:
			running = false
		case <-start:
			running = true
		case <-tick:
			if !running {
				break
			}
			pos += 2 * math.Pi * float64(*updateInterval) / float64(*period)
			update.Position.X = float32(*radius * math.Sin(pos))
			update.Position.Z = float32(*radius * math.Cos(pos))
			buf, err := proto.Marshal(update)
			if err != nil {
				log.Exit(err)
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, buf); err != nil {
				log.Exit(err)
			}
		}
	}
}

func api(ch chan<- struct{}, name, desc string) *pb.API {
	http.HandleFunc(fmt.Sprintf("/%s", name), func(w http.ResponseWriter, r *http.Request) {
		log.Infof("called %s", name)
		ch <- struct{}{}
	})
	return &pb.API{
		Endpoint:    fmt.Sprintf("http://%s/%s", *apiAddress, name),
		Name:        name,
		Description: desc,
	}
}
