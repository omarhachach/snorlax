package snorlax

import (
	"github.com/bwmarrin/discordgo"
)

// Module is used to import a modular package into the bot.
// This serves to make the bot modular and expandable.
type Module struct {
	Name       string
	Commands   map[string]Command
	RegisterOn map[string]interface{}
}

// Command holds the data and handler for a command.
type Command struct {
	Name       string
	Desc       string
	ModuleName string
	Handler    func(*discordgo.Session, *discordgo.MessageCreate)
}
