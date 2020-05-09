package server

import (
	"net"
	"time"
	"regexp"
)

const (
	// Status codes
	StatusHidden = 0
	StatusOnline = 1
	StatusAway   = 2
	StatusDND    = 3

	// User modes
	ModeRestricted = 0 // Do we need this?
	ModeUser = 1
	ModeAdmin = 2
	ModeOP = 99
)

var (
	// ServerUser = Client{Name: "System", LoginTime: time.Now(), Status: 1}

	reNewline = regexp.MustCompile(`\r?\n`)
)

// Server ...
type Server struct {
	StartTime time.Time
	Clients []*Client
	ServerUser *Client
	CommandHandler *CommandHandler
	// Server 
}

// Client ...
type Client struct {
	Socket net.Conn
	LoginTime time.Time
	Name string
	Status int
	Mode int // permissions
}

// CommandHandler ...
type CommandHandler struct {
	Commands map[string]*Command
	// Aliases map[string]*Command
}

// Command ...
type Command struct {
	Name string
	Description string
	Mode int
	Run CommandRunFunc
}

// CommandRunFunc ...
type CommandRunFunc func(*Server, *Client, []string)
