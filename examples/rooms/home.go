package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/dichro/intangible/server"

	pb "github.com/dichro/intangible"
	log "github.com/golang/glog"
)

var (
	httpRoot = flag.String("http_root", "", "directory containing static content to serve")
	port     = flag.Int("port", 8080, "listening port")
)

func main() {
	flag.Parse()
	room := server.NewRoom("home")
	static := []*pb.ServerUpdate{
		// Sculpture from the British Museum via Sketchfab, https://skfb.ly/BuOq
		{
			BoundingBox: &pb.Vector{3, 3, 3},
			Position:    &pb.Vector{-2, 1.5, 0},
			Rotation:    &pb.Vector{0, 90, 0},
			Rendering: &pb.Rendering{
				Obj: &pb.WavefrontOBJ{
					Obj: &pb.Source{Uri: "http://intangible-gallery.s3-website-us-west-1.amazonaws.com/sculpt.obj"},
				},
			},
		},
		{
			BoundingBox: &pb.Vector{3, 3, 3},
			Position:    &pb.Vector{2, 1.5, 0},
			Rendering: &pb.Rendering{
				Obj: &pb.WavefrontOBJ{
					Obj: &pb.Source{
						Uri:     "extincteur_obj.obj",
						Archive: &pb.Source{Uri: "http://www.oyonale.com/downloads/extincteur_obj.zip"},
					},
				},
			},
		},
		{
			BoundingBox: &pb.Vector{X: 10, Y: 1, Z: 10},
			Position:    &pb.Vector{Y: -0.5},
			Rendering: &pb.Rendering{
				Unit: &pb.UnitGeometry{
					Form:    pb.UnitGeometry_CUBE,
					Rescale: &pb.Vector{X: 10, Y: 1, Z: 10},
				},
			},
		},
		{
			BoundingBox: &pb.Vector{1, 3, 1},
			Position:    &pb.Vector{0, 1.5, 3},
			Api: []*pb.API{{
				// ws: endpoints are interpreted as connections to new rooms
				Endpoint:    "ws://intangible.gallery/portal/start",
				Name:        "teleport",
				Description: "teleport to gallery",
			}},
			// TODO(dichro): find a door rendering
		},
	}
	for _, obj := range static {
		room.Place(obj)
	}
	if len(*httpRoot) > 0 {
		http.Handle("/", http.FileServer(http.Dir(*httpRoot)))
	}
	http.Handle("/portal/start", server.NewWebsocketHandler(room))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
