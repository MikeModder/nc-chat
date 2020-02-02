package server

import (
	"fmt"
	"net"
	"time"
	"strings"
)

func NewServer() Server {
	return Server{}
}

func (s *Server) Run(address string, port int) {
	server, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(address), Port: port})
	if err != nil {
		panic(err)
	}

	for {
		nc, err := server.Accept()
		if err != nil || nc == nil {
			continue
		}

		s.AddClient(&Client{Socket: nc})
	}
}

func (s *Server) AddClient(c *Client) {
	c.LoginTime = time.Now()
	c.Name = "Unset" // TODO: set a a randomized name, to avoid confusion
	c.Status = 1

	s.SendToAll(fmt.Sprintf("%s joined\n", c.Name), &ServerUser)

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
			s.SendToAll(fmt.Sprintf("%s left (Reason: %s)\n", c.Name, reason), &ServerUser)

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
	switch (command) {
	case "nick":
		nick := reNewline.ReplaceAllString(args[0], "")
		if len(nick) > 24 {
			nick = nick[0:23]
		}

		if nick == "" {
			nick = "Unset"
		}

		if nick == "System" {
			nick = "Not System"
			invoker.Send("[System] Nice try, bud\n")
		}

		s.SendToAll(fmt.Sprintf("%s is now %s\n", invoker.Name, nick), &ServerUser)
		invoker.Name = nick
		break

	case "whisper":
		s.SendToUserByName(args[0], strings.Join(args[1:], " "), invoker)
		break

	default:
		invoker.Send("[System] Unknown command!\n")
	}
}

func (s *Server) SendToAll(message string, from *Client) {
	// fmt.Println(message)
	for _, c := range s.Clients {
		if c == from {
			continue
		}

		message = reNewline.ReplaceAllString(message, "")

		err := c.Send(fmt.Sprintf("[%s]: %s\n", from.Name, message))
		if err != nil {
			s.RemoveClient(c, "failed to send message to client")
		}
	}
}

func (s *Server) SendToUserByName(name string, message string, from *Client) {
	message = reNewline.ReplaceAllString(message, "")

	for _, c := range s.Clients {
		if c.Name == name {
			err := c.Send(fmt.Sprintf("[%s -> you]: %s\n", from.Name, message))
			if err != nil {
				from.Send("[System]: failed to send whisper...")
				s.RemoveClient(c, "failed to send whisper to client")
			}
			return
		}
	}

	from.Send("[System]: failed to send whisper, no user with that name\n")
}

// func (s *Server) FindUser(name string) Client {
// 	for _, c := range s.Clients {
// 		if c.Name == name {
// 			return *c
// 		}
// 	}

// 	return Client{}
// }