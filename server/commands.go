package server

import (
	"fmt"
	"hash/crc32"
	"strings"
	"time"
)

// RegisterCommands is called by the *server.Start() function to register commands
func RegisterCommands(c *CommandHandler) {
	c.AddCommand("nick", "Set your nickname", ModeUser, nickCmd)
	c.AddCommand("aboutme", "information about user", ModeUser, aboutmeCmd)
	c.AddCommand("whisper", "send a whisper/pm", ModeUser, whisperCmd)
	c.AddCommand("leave", "disconnect with a message", ModeUser, leaveCmd)
}

func nickCmd(s *Server, i *Client, args []string) {
	nick := reNewline.ReplaceAllString(strings.Join(args, " "), "")
	if len(nick) > 24 {
		i.SendSystemMessage("nicknames cannot be longer than 24 characters, yours has been truncated.")
		nick = nick[0:23]
	}

	if nick == "" {
		nick = fmt.Sprintf("%X", crc32.ChecksumIEEE([]byte(i.Socket.RemoteAddr().String())))
	}

	if nick == "System" {
		i.SendSystemMessage("invalid username!")
		return
	}

	// TODO: check if nickname is already in use
	s.Broadcast(i.Name + " is now " + nick)
	i.Name = nick
}

func aboutmeCmd(s *Server, i *Client, _ []string) {
	i.SendSystemMessage(fmt.Sprintf("You are %s logged in from %s for %s (mode: %d)", i.Name, "RemoteAddr()", time.Since(i.LoginTime).String(), i.Mode))
}

func whisperCmd(s *Server, i *Client, args []string) {
	// TODO: make this better
	s.SendToUserByName(args[0], strings.Join(args[1:], " "), i)
}

func leaveCmd(s *Server, i *Client, args []string) {
	msg := reNewline.ReplaceAllString(strings.Join(args, " "), "")
	if len(msg) > 64 {
		s.RemoveClient(i, msg[:63])
		return
	}

	s.RemoveClient(i, msg)
}
