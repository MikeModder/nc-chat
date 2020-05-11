package main

import (
	"flag"
	"fmt"

	"github.com/MikeModder/nc-chat/server"
)

func main() {
	fmt.Println("nc-chat server starting...")
	addr := flag.String("addr", "0.0.0.0", "Listen address for the server")
	port := flag.Int("port", 5006, "Port to run the server on")
	flag.Parse()

	fmt.Printf("running server on %s:%d\n", *addr, *port)
	app := server.NewServer()
	app.Run(*addr, *port)
}
