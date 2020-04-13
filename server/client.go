package server

import (
	"fmt"
	// "strings"
)

func (c *Client) SendSystemMessage(message string) (err error) {
	err = c.SendRaw(fmt.Sprintf("[Server]: %s\n", message))
	return
}

func (c *Client) SendMessageFromUser(message string, from *Client) (err error) {
	err = c.SendRaw(fmt.Sprintf("[%s]: %s\n", from.Name, message))
	return
}

// should I have a dedicated function to format messages befoer they get sent to the user?
// this would make consistency and updating formatting easier, as well as possibly allowing
// custom formatting via a configuration option.
func (c *Client) SendRaw(message string) (err error) {
	_, err = c.Socket.Write([]byte(message))
	return
// 	if err != nil {
// 		c.Kick(fmt.Sprintf("sever error: %v\n", err))
// 	}
}

func (c *Client) Kick(reason string) {
	//c.Send(fmt.Sprintf("You where kicked! Reason: %s\n", reason))
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