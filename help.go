package snorlax

import (
	"github.com/bwmarrin/discordgo"
)

var helpMessage string

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

	helpModule := &Module{
		Name:     moduleName,
		Commands: commands,
		Init:     helpInit,
	}

	internalModules[helpModule.Name] = helpModule
}

// moduleCommands holds a list of modules, with their respective commands.
var moduleCommands = map[string]map[string]*Command{}

func helpInit(s *Snorlax) {
	for _, module := range s.Modules {
		moduleCommands[module.Name] = module.Commands
	}
}

func helpHandler(s *Snorlax, m *discordgo.MessageCreate) {

}
