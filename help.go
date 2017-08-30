package snorlax

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	helpCommand := &Command{
		Command: ".help",
		Alias:   ".h",
		Desc:    "Help shows you a help menu for a given module, or a list of modules.",
		Usage:   ".help [module-name]",
		Init:    helpInit,
		Handler: helpHandler,
	}

	reservedCommands[helpCommand.Command] = helpCommand
}

func helpInit(s *Snorlax) {

}

func helpHandler(s *Snorlax, m *discordgo.MessageCreate) {

}
