package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"

	pb "github.com/dichro/intangible"
	avro "github.com/elodina/go-avro"
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
		center  = pb.Vector{Y: 1}
		tick    = time.Tick(*updateInterval)
		pos     = 0.0
		running = true

		// async notification channels from incoming API calls
		stop  = make(chan struct{})
		start = make(chan struct{})
		orbit = make(chan *pb.Vector)

		// keep a single Object entry around and just update it, rather than re-creating
		update = &pb.Object{
			Id: "self",
			Api: []*pb.API{
				simpleAPI(stop, "stop", "stop orbiter"),
				simpleAPI(start, "start", "start orbiter"),
				vectorAPI(orbit, "orbit", "set orbit center"),
			},
			Position: &pb.Vector{},
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
	)
	go http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	for {
		select {
		case <-stop:
			running = false
		case <-start:
			running = true
		case v := <-orbit:
			center = *v
		case <-tick:
			if !running {
				break
			}
			pos += 2 * math.Pi * float64(*updateInterval) / float64(*period)
			update.Position.X = center.X + float32(*radius*math.Sin(pos))
			update.Position.Y = center.Y
			update.Position.Z = center.Z + float32(*radius*math.Cos(pos))
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

// TODO(dichro): there should be a way to generate the schema from an annotated struct. Furthermore, there should probably be a predefined well-known-schema for things like this.
const requestSchema = `{"namespace": "intangible.orbiter",
			 "type": "record",
			 "name": "orbit",
			 "fields": [
			     {"name": "x", "type": "float"},
			     {"name": "y", "type": "float"},
			     {"name": "z", "type": "float"}
			 ]
			}`

var schema = avro.MustParseSchema(requestSchema)

func vectorAPI(ch chan<- *pb.Vector, name, desc string) *pb.API {
	reader := avro.NewGenericDatumReader()
	reader.SetSchema(schema)
	http.HandleFunc(fmt.Sprintf("/%s", name), func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			return
		}
		msg := avro.NewGenericRecord(schema)
		decoder := avro.NewBinaryDecoder(buf)
		if err := reader.Read(msg, decoder); err != nil {
			log.Error(err)
			return
		}
		ch <- &pb.Vector{
			X: msg.Get("x").(float32),
			Y: msg.Get("y").(float32),
			Z: msg.Get("z").(float32),
		}
	})
	api := api(name, desc)
	api.AvroRequestSchema = requestSchema
	return api
}

func simpleAPI(ch chan<- struct{}, name, desc string) *pb.API {
	http.HandleFunc(fmt.Sprintf("/%s", name), func(w http.ResponseWriter, r *http.Request) {
		log.Infof("called %s", name)
		ch <- struct{}{}
	})
	return api(name, desc)
}

func api(name, desc string) *pb.API {
	return &pb.API{
		Endpoint:    fmt.Sprintf("http://%s/%s", *apiAddress, name),
		Name:        name,
		Description: desc,
	}
}
