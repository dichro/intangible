package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/dichro/intangible/async"
	"github.com/dichro/intangible/server"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"

	pb "github.com/dichro/intangible"
	log "github.com/golang/glog"
)

var (
	httpRoot = flag.String("http_root", "", "directory containing static content to serve")
	port     = flag.Int("port", 8080, "listening port")
)

func main() {
	flag.Parse()
	room := async.NewLatestSnapshot(20)
	static := []*pb.Object{
		{
			Id:          "sculpture",
			BoundingBox: &pb.Vector{3, 3, 3},
			Position:    &pb.Vector{-2, 1.5, 0},
			Rotation:    &pb.Vector{0, 90, 0},
			Rendering: &pb.Rendering{
				Mesh: &pb.Mesh{
					SourceUri: "http://intangible-gallery.s3-website-us-west-1.amazonaws.com/sculpt.obj",
					Rotation:  &pb.Vector{-90, 53, 0},
				},
				Texture: &pb.Texture{
					SourceUri: "http://intangible-gallery.s3-website-us-west-1.amazonaws.com/tex_0.jpg",
				},
			},
		},
		{
			Id:          "floor",
			BoundingBox: &pb.Vector{X: 10, Y: 1, Z: 10},
			Position:    &pb.Vector{Y: -0.5},
			Rendering: &pb.Rendering{
				Mesh: &pb.Mesh{
					SourceUri: "platonic:cube",
					Rescale:   &pb.Vector{X: 10, Y: 1, Z: 10},
				},
				Texture: &pb.Texture{
					SourceUri: "colour:grey",
				},
			},
		},
		{
			Id:          "hoa-hakananaia",
			BoundingBox: &pb.Vector{X: 3, Y: 3, Z: 3},
			Position:    &pb.Vector{Y: 1.5, X: 2},
			Rotation:    &pb.Vector{Y: 290},
			Rendering: &pb.Rendering{
				Mesh: &pb.Mesh{
					SourceUri: "http://intangible-gallery.s3-website-us-west-1.amazonaws.com/hoa-hakananaia/sculpt.obj",
					Rotation:  &pb.Vector{X: 90},
				},
				Texture: &pb.Texture{
					SourceUri: "http://intangible-gallery.s3-website-us-west-1.amazonaws.com/hoa-hakananaia/tex_0.jpg",
				},
			},
		},
		{
			Id:          "door",
			BoundingBox: &pb.Vector{1, 3, 1},
			Position:    &pb.Vector{0, 1.5, 3},
			Api: []*pb.API{{
				Endpoint:    "ws://ws.intangible.gallery/",
				Name:        "teleport",
				Description: "teleport to gallery",
			}},
		},
	}
	for _, obj := range static {
		buf, err := proto.Marshal(obj)
		if err != nil {
			log.Fatal(err)
		}
		room.Update(obj.Id, 1, buf)
	}
	if len(*httpRoot) > 0 {
		http.Handle("/", http.FileServer(http.Dir(*httpRoot)))
	}
	http.Handle("/ws", prometheus.InstrumentHandler("websockets", server.NewHandler(room)))
	http.Handle("/metrics", prometheus.Handler())
	prometheus.MustRegister(prometheus.NewProcessCollector(os.Getpid(), "server"))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
