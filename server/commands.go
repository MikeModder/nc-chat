package server

import (
	"fmt"
	"hash/crc32"
	"strings"
	"time"
)

func RegisterCommands(c *CommandHandler) {
	//c.AddCommand("nick", "Set your nickname", ModeUser, nickCmd)
	c.AddCommand("aboutme", "information about user", ModeUser, aboutmeCmd)
	c.AddCommand("whisper", "send a whisper/pm", ModeUser, whisperCmd)
	c.AddCommand("leave", "disconnect with a message", ModeUser, leaveCmd)

	// Authentication shit
	c.AddCommand("login", "authenticate yourself", ModeUnauthenticated, loginCmd)
	c.AddCommand("register", "sign up, I guess", ModeUnauthenticated, registerCmd)
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

func registerCmd(s *Server, i *Client, args []string) {
	// Register will create a new user with the given username and password (assuming it doesn't
	// already exist) and then overwrite the current user's details. This needs to check both 
	// authentication status and if the username is claimed.

	// First, check if they're already logged in
	if i.Mode >= ModeUser {
		i.SendSystemMessage("You are already registered and logged in!")
		return
	}

	// Secondly, check amount of arguments
	if len(args) < 3 {
		i.SendSystemMessage("You didn't pass in enough arguments!")
		i.SendSystemMessage("Required arguments: <username> <password> <password>")
		return
	}

	args[2] = reNewline.ReplaceAllString(args[2], "")
	// Third, make sure the passwords match
	fmt.Printf("%v\n", args)
	if args[1] != args[2] {
		i.SendSystemMessage("Passwords don't match!")
		return
	}

	// Probably should have done this earlier, but check if there's a user with that name already
	// TODO: check against database using wrapper function (UserExists()?)

	// Create user in database
	err := CreateUser(args[0], args[1])
	if err != nil {
		// TODO: log errors in some sensible manner, even if it's just the console
		i.SendSystemMessage("Your registration couldn't be processed at this time, please try again later!")
		fmt.Printf("[error] failed to CreateUser(): %v\n", err)
		return
	}
	// Overwrite current user info with chosen details
	i.Name = args[0]
	i.Mode = ModeUser

	// Inform user of success
	i.SendSystemMessage(fmt.Sprintf("You have been registered and logged in! In the future you can log in with /login %s %s", args[0], args[1]))
}

func loginCmd(s *Server, i *Client, args []string) {
	// Make sure the user isn't already logged in
	if i.Mode >= ModeUser {
		i.SendSystemMessage("You are already logged in!")
	}

	if len(args) < 2 {
		// TODO: more helpful error message
		i.SendSystemMessage("You need to pass in two arguments")
		return
	}

	u, err := GetUserByName(args[0])
	if err != nil {
		fmt.Printf("[error] failed to log in: %v\n", err)
		i.SendSystemMessage("failed to log in")
		return
	}

	args[1] = reNewline.ReplaceAllString(args[1], "")
	err = CheckPassword(args[1], u.Password)
	if err != nil {
		i.SendSystemMessage("failed to log in, is your password incorrect?")
		return
	}

	// Their password matched, I guess. Let 'em in!
	oldname := i.Name
	i.Name = u.Name
	i.Mode = u.Mode

	// Notify the user and everyone else they authenticated
	i.SendSystemMessage(fmt.Sprintf("Logged in as %s!", u.Name))
	s.Broadcast(fmt.Sprintf("%s authenticated as %s!", oldname, u.Name))
}
