package main

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
	BotID string
	TelnetHost string
	TelnetPassword string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&TelnetHost, "h", "localhost:8081", "Telnet Server")
	flag.StringVar(&TelnetPassword, "p", "", "Telnet Password")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	// Handle all commands with leading exclamation point
	switch m.Content {
	case "!help":
		_, _ = s.ChannelMessageSend(m.ChannelID, "This bot has the following commands enabled:\n!status - server status\n!players - list all players connected\n!time - display the current day and time of the server\n!version - print the version of the server")
		return
	case "!status":
		_, _ = s.ChannelMessageSend(m.ChannelID, getServerStatus())
		return
	case "!players":
		_, _ = s.ChannelMessageSend(m.ChannelID, getPlayerList())
		return
	case "!time":
		_, _ = s.ChannelMessageSend(m.ChannelID, getServerTime())
		return
	case "!version":
		_, _ = s.ChannelMessageSend(m.ChannelID, getServerVersion())
		return
	}

	// Else, relay say message to game server
	sendServerSay(m)
}

func sendServerSay(m *discordgo.MessageCreate) {
	//Send message to game server over telnet connection
	sendTelnetMessage("say \"" + m.Author.Username + " (Discord): " + m.Content + "\"")
}

func getServerStatus() string {
	//Get player list from telnet connection
	_, err := sendTelnetMessage("version")
	if err == nil {
		return "Game server is running"
	} else {
		return "Error contacting game server"
	}
}

func getPlayerList() string {
	//Get player list from telnet connection
	result, err := sendTelnetMessage("lp")
	if err == nil {
		return result
	} else {
		return "Error contacting game server"
	}
}

func getServerTime() string {
	//Get server time from telnet connection
	result, err := sendTelnetMessage("gt")
	if err == nil {
		return result
	} else {
		return "Error contacting game server"
	}
}

func getServerVersion() string {
	//Get server version from telnet connection
	result, err := sendTelnetMessage("version")
	if err == nil {
		return result
	} else {
		return "Error contacting game server"
	}
}
