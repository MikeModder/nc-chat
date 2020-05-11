package server

import (
	"fmt"
	// "strings"
)

// SendSystemMessage sends a message (from the system) to the client
func (c *Client) SendSystemMessage(message string) (err error) {
	err = c.SendRaw(fmt.Sprintf("[Server]: %s\n", message))
	return
}

// SendMessageFromUser sends a message to the client, from another client
func (c *Client) SendMessageFromUser(message string, from *Client) (err error) {
	err = c.SendRaw(fmt.Sprintf("[%s]: %s\n", from.Name, message))
	return
}

// should I have a dedicated function to format messages befoer they get sent to the user?
// this would make consistency and updating formatting easier, as well as possibly allowing
// custom formatting via a configuration option.
// SendRaw sends a raw string to the client.
func (c *Client) SendRaw(message string) (err error) {
	_, err = c.Socket.Write([]byte(message))
	return
// 	if err != nil {
// 		c.Kick(fmt.Sprintf("sever error: %v\n", err))
// 	}
}

// Kick sends a kick message to the client and closes the socket
func (c *Client) Kick(reason string) {
	//c.Send(fmt.Sprintf("You where kicked! Reason: %s\n", reason))
	c.SendSystemMessage("You where kicked! Reason: " + reason)
	c.Socket.Close()
}

// func (c *Client) Leave(message string) {

// }
