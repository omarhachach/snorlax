package snorlax

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	// Commands is a list of all the registered commands, and their assocated Command type.
	Commands map[string]*Command
)

// Snorlax is the bot type.
type Snorlax struct {
	Discord *discordgo.Session
	Modules map[string]Module
	Log     *logrus.Logger
}

// NewBot returns a new bot type.
func NewBot(discord *discordgo.Session) *Snorlax {
	Commands = make(map[string]*Command)

	return &Snorlax{
		Discord: discord,
		Modules: make(map[string]Module),
		Log:     logrus.New(),
	}
}

// RegisterModule allows you to register a module to expand the bot.
func (s *Snorlax) RegisterModule(module Module) {
	_, moduleExist := s.Modules[module.Name]
	if moduleExist {
		s.Log.Error("Failed to load module: " + module.Name + ".\nModule with same name has already been registered.")
		return
	}

	for _, command := range module.Commands {
		existingCommand, commandExist := Commands[command.Name]
		if commandExist {
			s.Log.Error("Failed to load module: " + module.Name +
				".\nModule " + existingCommand.ModuleName + "has already registered command/alias: " + command.Name)
			return
		}

		if command.Alias != "" {
			existingAlias, aliasExist := Commands[command.Alias]
			if aliasExist {
				s.Log.Error("Failed to load module: " + module.Name +
					".\nModule " + existingAlias.ModuleName + "has already registered command/alias: " + command.Alias)
				return
			}

			s.Log.Debug("Registered Alias: " + command.Alias)
			Commands[command.Alias] = command
		}

		s.Log.Debug("Registered Command: " + command.Name)
		Commands[command.Name] = command
	}

	s.Modules[module.Name] = module
	s.Log.Info("Loaded module: " + module.Name)
}

// Start opens a connection to Discord, and initiliazes the bot.
func (s *Snorlax) Start() {
	s.Discord.AddHandler(onMessageCreate(s))

	s.Log.SetLevel(logrus.DebugLevel)
	err := s.Discord.Open()
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error establishing connection to Discord.")
		return
	}

	s.Log.Info("Snorlax has been woken!")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	s.Log.Info("Snorlax is now sleeping.")
	err = s.Discord.Close()
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error closing Discord session.")
	}
}
