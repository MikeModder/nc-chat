package server

import (
	"fmt"
	// "strings"
)

func (c *Client) Send(message string) (err error) {
	_, err = c.Socket.Write([]byte(message))
	return
// 	if err != nil {
// 		c.Kick(fmt.Sprintf("sever error: %v\n", err))
// 	}
}

func (c *Client) Kick(reason string) {
	c.Send(fmt.Sprintf("You where kicked! Reason: %s\n", reason))
	c.Socket.Close()
}

// func (c *Client) Leave(message string) {

// }

// func (c *Client) HandleCommand(command string, args []string) {
// 	switch (command) {
// 	case "nick":
// 		nick := strings.Join(args, " ")
// 		if len(nick) > 24 {
// 			nick = nick[0:23]
// 		}
		
// 		break
// 	}
// }