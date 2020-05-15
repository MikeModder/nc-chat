package server

import (
	"fmt"
	"net"
	"time"
	"strings"
	"hash/crc32"
)

// NewServer creates a new server with a default command handler and System user
func NewServer() Server {
	return Server{
		ServerUser: &Client{
			Name: "System",
			LoginTime: time.Now(),
			Status: 1,
			Mode: ModeAdmin,
		},
		CommandHandler: NewCommandHandler(),
	}
}

// Run starts up the server, and handles new connections
func (s *Server) Run(address string, port int) {
	server, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(address), Port: port})
	if err != nil {
		panic(err)
	}

	// Register commands into command handler
	RegisterCommands(s.CommandHandler)
	// Initialize database
	InitDB()

	for {
		nc, err := server.Accept()
		if err != nil || nc == nil {
			continue
		}

		// s.AddClient(&Client{Socket: nc})
		s.AddClient(&Client{
			Socket: nc,
			LoginTime: time.Now(),
			Name: fmt.Sprintf("%X", crc32.ChecksumIEEE([]byte(nc.RemoteAddr().String()))), // TODO: how heavy is this? it's only run once on first connection
			Status: StatusOnline,
			Mode: ModeUnauthenticated,
		})
	}
}

// AddClient adds the client to the master list and starts a goroutine to handle their socket
func (s *Server) AddClient(c *Client) {
	c.SendSystemMessage("Welcome to nc-chat server!")
	c.SendSystemMessage("Be sure to login or register with /login or /register respectively!")
	s.Broadcast(c.Name + " joined")

	s.Clients = append(s.Clients, c)
	go s.HandleClient(c)
}

// RemoveClient closes the passed Client's socket and notifies the other users, optionally with a message
func (s *Server) RemoveClient(c *Client, reason string) {
	c.Socket.Close()

	if reason == "" { reason = "None" }

	var i int
	for i = range s.Clients {
		if s.Clients[i] == c {
			s.Clients = append(s.Clients[:i], s.Clients[i+1:]...)
			s.Broadcast(fmt.Sprintf("%s left (Reason: %s)", c.Name, reason))

			break
		}
	}

}

// HandleClient is run as a goroutine for each connected client. It handles reading from the socket and determining if a command needs to be run
func (s *Server) HandleClient(c *Client) {
	buf := make([]byte, 2048)

	for {
		read, err := c.Socket.Read(buf)
		if err != nil {
			s.RemoveClient(c, "failed to read from socket")
			break
		}

		str := string(buf[:read])

		if reNewline.ReplaceAllString(str, "") == "" { continue }

		if strings.HasPrefix(str, "/") {
			// Get ready to run a command
			split := strings.Split(str, " ")
			s.HandleCommand(c, strings.TrimPrefix(split[0], "/"), split[1:])
			continue
		}

		// Do this check after the command handler, since the user has to be able to log in
		if c.Mode < ModeUser {
			c.SendSystemMessage("You are either not logged in or have been restricted!")
			continue
		}
		s.SendToAll(str, c)
	}
}

// HandleCommand doesn't directly handle the command, rather it strips newlines and shows a neat error message
func (s *Server) HandleCommand(invoker *Client, command string, args []string) {
	command = reNewline.ReplaceAllString(command, "")

	// Strupt newlines from args as well
	for i := 0; i < len(args)-1; i++ {
		args[i] = reNewline.ReplaceAllString(args[i], "")
	}

	ok, err := s.CommandHandler.ExecuteCommand(s, invoker, command, args)
	if !ok {
		invoker.SendSystemMessage("comand returned error: " + err)
	}
}

// Broadcast sends a (system) message to all connected clients
func (s *Server) Broadcast(message string) {
	for i := len(s.Clients)-1; i >= 0; i-- {
		s.Clients[i].SendSystemMessage(message)
	}
}

// SendToAll sends a message to all clients, originating from another client
func (s *Server) SendToAll(message string, from *Client) {
	for i := len(s.Clients)-1; i >= 0; i-- {
		// Don't forward messages to restricted or unauthorized users
		if s.Clients[i].Mode < ModeUser { continue; }
		if s.Clients[i] == from {
			continue
		}

		message = reNewline.ReplaceAllString(message, "")

		err := s.Clients[i].SendMessageFromUser(message, from)
		if err != nil {
			s.RemoveClient(s.Clients[i], "failed to send message to client")
		}
	}
}

// TODO: SendToUserByName assumes it's being invoked with the whisper command
// a dedicated function for whispering would probably be better
// make the whisper function a wrapper on top of this? but that's 3(?) layers of abstraction already
// SendToUserByName sends a whisper to a user with a given username, from another client
func (s *Server) SendToUserByName(name string, message string, from *Client) {
	message = reNewline.ReplaceAllString(message, "")

	for i := len(s.Clients)-1; i >= 0; i-- {
		if s.Clients[i].Name == name {
			// Don't allow sending whispers to restricted/unauthorized accounts
			if s.Clients[i].Mode < ModeUser {
				from.SendSystemMessage("That user is either restricted or not signed in, so you cannot send them a whisper.")
				return
			}

			err := s.Clients[i].SendRaw(fmt.Sprintf("[%s -> you]: %s\n", from.Name, message)) // TODO: revise this to use a wrapper function
			// err := c.SendMessageFromUser(message, from)
			if err != nil {
				from.SendSystemMessage("Failed to send whisper")
				// s.RemoveClient(c, "failed to send whisper to client")
			}
			return
		}
	}

	from.SendSystemMessage("Failed to send whisper to user. Is there a user connected with that username?")
}

// func (s *Server) FindUser(name string) Client {
// 	for _, c := range s.Clients {
// 		if c.Name == name {
// 			return *c
// 		}
// 	}

// 	return Client{}
// }
