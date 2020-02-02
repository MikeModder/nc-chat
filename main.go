package main

import (
	"github.com/MikeModder/nc-chat/server"
)

func main() {
	app := server.NewServer()
	app.Run("0.0.0.0", 5006)
}