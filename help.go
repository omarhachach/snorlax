package snorlax

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	moduleName := "Help"
	commands := map[string]*Command{}

	helpCommand := &Command{
		Command: ".help",
		Alias:   ".h",
		Desc:    "Help shows you a help menu for a given module, or a list of modules.",
		Usage:   ".help [module-name]",
		Handler: helpHandler,
	}

	commands[helpCommand.Command] = helpCommand

	helpModule := &internalModule{
		Name:     moduleName,
		Commands: commands,
		Init:     helpInit,
	}

	internalModules[helpModule.Name] = helpModule
}

func helpInit(s *Snorlax) {

}

func helpHandler(s *Snorlax, m *discordgo.MessageCreate) {

}
