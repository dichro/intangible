This is a simple implementation of a VR room. To start an instance up, say:

```
go get -u github.com/dichro/intangible/examples/rooms
go install github.com/dichro/intangible/examples/rooms
$GOPATH/bin/rooms --alsologtostderr
```

This will start up the room server on port 8080.

To connect to it, get the Tesseract client for the Vive (Windows) from [itch.io](https://intangible.itch.io/tesseract).

Running the client will by default connect to the public intangible.gallery room. To connect to your freshly installed room instead, add the command-line flags `+startRoom ws://localhost:8080/portal/start`