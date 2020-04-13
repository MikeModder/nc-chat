package server

import (
	"fmt"
	"net"
	"time"
	"strings"
	"hash/crc32"
)

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

func (s *Server) Run(address string, port int) {
	server, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(address), Port: port})
	if err != nil {
		panic(err)
	}

	// Register commands into command handler
	RegisterCommands(s.CommandHandler)

	for {
		nc, err := server.Accept()
		if err != nil || nc == nil {
			continue
		}

		

		// s.AddClient(&Client{Socket: nc})
		s.AddClient(&Client{
			Socket: nc,
			LoginTime: time.Now(),
			Name: fmt.Sprintf("%X", crc32.ChecksumIEEE([]byte(nc.RemoteAddr().String()))), // TODO: how heavy is this? it's on
			Status: StatusOnline,
			Mode: ModeUser,
		})
	}
}

func (s *Server) AddClient(c *Client) {
	// c.LoginTime = time.Now()
	// c.Name = "Unset" // TODO: set a a randomized name, to avoid confusion
	// c.Status = 1

	c.SendSystemMessage("Welcome to nc-chat server!")
	s.Broadcast(c.Name + " joined")

	s.Clients = append(s.Clients, c)
	go s.HandleClient(c)
}

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
			// buf = make([]byte, 2048)
			continue
		}

		s.SendToAll(str, c)
	}
}

func (s *Server) HandleCommand(invoker *Client, command string, args []string) {
	command = reNewline.ReplaceAllString(command, "")
	// fmt.Println(command)

	ok, err := s.CommandHandler.ExecuteCommand(s, invoker, command, args)
	if !ok {
		invoker.SendSystemMessage("comand returned error: " + err)
	}
}

func (s *Server) Broadcast(message string) {
	for i := len(s.Clients)-1; i >= 0; i-- {
		s.Clients[i].SendSystemMessage(message)
	}
}

func (s *Server) SendToAll(message string, from *Client) {
	// fmt.Println(message)
	for _, c := range s.Clients {
		if c == from {
			continue
		}

		message = reNewline.ReplaceAllString(message, "")

		err := c.SendMessageFromUser(message, from)
		// err := c.Send(fmt.Sprintf("[%s]: %s\n", from.Name, message))
		if err != nil {
			s.RemoveClient(c, "failed to send message to client")
		}
	}
}

// TODO: SendToUserByName assumes it's being invoked with the whisper command
// a dedicated function for whispering would probably be better
// make the whisper function a wrapper on top of this? but that's 3(?) layers of abstraction already
func (s *Server) SendToUserByName(name string, message string, from *Client) {
	message = reNewline.ReplaceAllString(message, "")

	for _, c := range s.Clients {
		if c.Name == name {
			err := c.SendRaw(fmt.Sprintf("[%s -> you]: %s\n", from.Name, message)) // TODO: revise this to use a wrapper function
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