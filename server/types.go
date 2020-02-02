package server

import (
	"net"
	"time"
	"regexp"
)

const (
	StatusHidden = 0
	StatusOnline = 1
	StatusAway   = 2
	StatusDND    = 3
)

var (
	ServerUser = Client{Name: "System", LoginTime: time.Now(), Status: 1}

	reNewline = regexp.MustCompile(`\r?\n`)
)

type Server struct {
	Clients []*Client
	// Server 
}


type Client struct {
	Socket net.Conn
	LoginTime time.Time
	Name string
	Status int
}