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
	}
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

		

		// s.AddClient(&Client{Socket: nc})
		s.AddClient(&Client{
			Socket: nc,
			LoginTime: time.Now(),
			Name: fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(nc.RemoteAddr().String()))), // TODO: how heavy is this? it's on
			Status: StatusOnline,
		})
	}
}

func (s *Server) AddClient(c *Client) {
	// c.LoginTime = time.Now()
	// c.Name = "Unset" // TODO: set a a randomized name, to avoid confusion
	// c.Status = 1

	c.SendSystemMessage("Welcome to nc-chat server!")
	s.SendToAll(fmt.Sprintf("%s joined\n", c.Name), s.ServerUser)

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
			s.SendToAll(fmt.Sprintf("%s left (Reason: %s)\n", c.Name, reason), s.ServerUser)

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
			invoker.SendSystemMessage("Nice try, bud")
			//invoker.Send("[System] Nice try, bud\n")
		}

		s.SendToAll(fmt.Sprintf("%s is now %s", invoker.Name, nick), s.ServerUser)
		invoker.Name = nick
		break
	case "aboutme":
		invoker.SendMessageFromUser(fmt.Sprintf("You are %s logged in from %s for %s", invoker.Name, invoker.Socket.RemoteAddr(), time.Since(invoker.LoginTime).String()), s.ServerUser)
		break

	case "whisper":
		s.SendToUserByName(args[0], strings.Join(args[1:], " "), invoker)
		break

	default:
		// invoker.Send("[System] Unknown command!\n")
		invoker.SendSystemMessage("Unkown command!")
		break
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

func (s *Server) SendToUserByName(name string, message string, from *Client) {
	message = reNewline.ReplaceAllString(message, "")

	for _, c := range s.Clients {
		if c.Name == name {
			err := c.SendRaw(fmt.Sprintf("[%s -> you]: %s\n", from.Name, message)) // TODO: revise this to use a wrapper function
			// err := c.SendMessageFromUser(message, from)
			if err != nil {
				from.SendSystemMessage("Failed to send whisper")
				s.RemoveClient(c, "failed to send whisper to client")
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