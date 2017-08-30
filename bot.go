package snorlax

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var reservedCommands = map[string]*Command{}

// Snorlax is the bot type.
type Snorlax struct {
	Commands map[string]*Command
	Modules  map[string]*Module
	Session  *discordgo.Session
	Log      *logrus.Logger
	token    string
	config   *Config
}

// Config holds the options for the bot.
type Config struct {
	Debug bool
}

// New returns a new bot type.
func New(token string, config *Config) *Snorlax {
	return &Snorlax{
		Commands: map[string]*Command{},
		Modules:  map[string]*Module{},
		Log:      logrus.New(),
		token:    token,
		config:   config,
	}
}

// RegisterModule allows you to register a module to expand the bot.
func (s *Snorlax) RegisterModule(module *Module) {
	_, moduleExist := s.Modules[module.Name]
	if moduleExist {
		s.Log.Error("Failed to load module: " + module.Name + ".\nModule with same name has already been registered.")
		return
	}

	for _, command := range module.Commands {
		existingCommand, commandExist := s.Commands[command.Command]
		if commandExist {
			s.Log.Error("Failed to load module: " + module.Name +
				".\nModule " + existingCommand.ModuleName + "has already registered command/alias: " + command.Command)
			return
		}

		if command.Alias != "" {
			existingAlias, aliasExist := s.Commands[command.Alias]
			if aliasExist {
				s.Log.Error("Failed to load module: " + module.Name +
					".\nModule " + existingAlias.ModuleName + "has already registered command/alias: " + command.Alias)
				return
			}

			s.Log.Debug("Registered Alias: " + command.Alias)
			s.Commands[command.Alias] = command
		}

		s.Log.Debug("Registered Command: " + command.Command)
		s.Commands[command.Command] = command
	}

	s.Modules[module.Name] = module
	s.Log.Info("Loaded module: " + module.Name)
}

// Start opens a connection to Discord, and initiliazes the bot.
func (s *Snorlax) Start() {
	discord, err := discordgo.New("Bot " + s.token)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to create the Discord session")
		return
	}
	s.Session = discord

	s.Session.AddHandler(onMessageCreate(s))

	if s.config.Debug {
		s.Log.SetLevel(logrus.DebugLevel)
	}
	err = s.Session.Open()
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
	err = s.Session.Close()
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error closing Discord session.")
	}
}
