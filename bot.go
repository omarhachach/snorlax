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
	config  *Config
}

// Config holds the options for the bot.
type Config struct {
	Debug bool
}

// NewBot returns a new bot type.
func NewBot(discord *discordgo.Session, config *Config) *Snorlax {
	Commands = make(map[string]*Command)

	return &Snorlax{
		Discord: discord,
		Modules: make(map[string]Module),
		Log:     logrus.New(),
		config:  config,
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
		existingCommand, commandExist := Commands[command.Command]
		if commandExist {
			s.Log.Error("Failed to load module: " + module.Name +
				".\nModule " + existingCommand.ModuleName + "has already registered command/alias: " + command.Command)
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

		s.Log.Debug("Registered Command: " + command.Command)
		Commands[command.Command] = command
	}

	s.Modules[module.Name] = module
	s.Log.Info("Loaded module: " + module.Name)
}

// Start opens a connection to Discord, and initiliazes the bot.
func (s *Snorlax) Start() {
	s.Discord.AddHandler(onMessageCreate(s))

	if s.config.Debug {
		s.Log.SetLevel(logrus.DebugLevel)
	}
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
