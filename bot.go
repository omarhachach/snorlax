package snorlax

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	// Commands is a list of all the registered commands, and their assocated Command type.
	Commands map[string]*Command
)

// Snorlax is the bot type.
type Snorlax struct {
	Discord *discordgo.Session
	Modules map[string]Module
}

// NewBot returns a new bot type.
func NewBot(discord *discordgo.Session) *Snorlax {
	Commands = make(map[string]*Command)

	return &Snorlax{
		Discord: discord,
		Modules: make(map[string]Module),
	}
}

// RegisterModule allows you to register a module to expand the bot.
func (s *Snorlax) RegisterModule(module Module) {
	_, moduleExist := s.Modules[module.Name]
	if moduleExist {
		log.Error("Failed to load module: " + module.Name + ".\nModule with same name has already been registered.")
		return
	}

	for _, command := range module.Commands {

		existingCommand, commandExist := Commands[command.Name]
		if commandExist {
			log.Error("Failed to load module: " + module.Name +
				".\nModule " + existingCommand.ModuleName + "has already registered command/alias: " + command.Name)
			return
		}

		if command.Alias != "" {
			existingAlias, aliasExist := Commands[command.Alias]
			if aliasExist {
				log.Error("Failed to load module: " + module.Name +
					".\nModule " + existingAlias.ModuleName + "has already registered command/alias: " + command.Name)
				return
			}

			Commands[command.Alias] = command
		}

		Commands[command.Name] = command
	}

	s.Modules[module.Name] = module
	log.Info("Loaded module: " + module.Name)
}

// Start opens a connection to Discord, and initiliazes the bot.
func (s *Snorlax) Start() {
	s.Discord.AddHandler(onMessageCreate)

	err := s.Discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error establishing connection to Discord.")
		return
	}

	log.Info("Snorlax has been woken!")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	log.Info("Snorlax is now sleeping.")
	err = s.Discord.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error closing Discord session.")
	}
}
