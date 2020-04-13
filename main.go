package main

import (
	"flag"

	"github.com/MikeModder/nc-chat/server"
)

func main() {
	addr := flag.String("addr", "0.0.0.0", "Listen address for the server")
	port := flag.Int("port", 5006, "Port to run the server on")
	flag.Parse()

	app := server.NewServer()
	app.Run(*addr, *port)
}
