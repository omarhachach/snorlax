package snorlax

import (
	"database/sql"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax/models"
	"github.com/sirupsen/logrus"
)

var internalModules = map[string]*Module{}

// Version is the Go version.
const Version = "0.1.0"

// Snorlax is the bot type.
type Snorlax struct {
	Commands map[string]*Command
	Modules  map[string]*Module
	Session  *discordgo.Session
	Log      *logrus.Logger
	Mutex    *sync.Mutex
	DB       *sql.DB
	config   *Config
}

// New returns a new bot type.
func New(config *Config) *Snorlax {
	s := &Snorlax{
		Commands: map[string]*Command{},
		Modules:  map[string]*Module{},
		Log:      logrus.New(),
		Mutex:    &sync.Mutex{},
		config:   config,
	}

	if s.config.Debug {
		s.Log.SetLevel(logrus.DebugLevel)
	}

	for _, internalModule := range internalModules {
		s.RegisterModule(internalModule)
	}

	return s
}

// RegisterModule allows you to register a module to expand the bot.
func (s *Snorlax) RegisterModule(module *Module) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	_, moduleExist := s.Modules[module.Name]
	if moduleExist {
		s.Log.Error("Failed to load module: " + module.Name + ".\nModule with same name has already been registered.")
		return
	}

	for commandName, command := range module.Commands {
		existingCommand, commandExist := s.Commands[command.Command]
		if commandExist {
			s.Log.Error("Failed to load module: " + module.Name +
				".\nModule " + existingCommand.ModuleName + "has already registered command/alias: " + commandName)
			return
		}

		s.Log.Debug("Registered Command: " + commandName)
		s.Commands[commandName] = command

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
	}

	s.Modules[module.Name] = module
	s.Log.Info("Loaded module: " + module.Name)
}

// RegisterModules registers a list of modules.
func (s *Snorlax) RegisterModules(modules ...*Module) {
	for _, module := range modules {
		s.RegisterModule(module)
	}
}

// Start opens a connection to Discord, and initiliazes the bot.
func (s *Snorlax) Start() {
	go func() {
		s.Mutex.Lock()
		for _, module := range s.Modules {
			if module.Init != nil {
				go module.Init(s)
			}
		}
		s.Mutex.Unlock()
	}()

	discord, err := discordgo.New("Bot " + s.config.Token)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to create the Discord session")
		return
	}
	s.Session = discord

	s.Session.AddHandler(onMessageCreate(s))

	err = s.Session.Open()
	if err != nil {
		s.Log.WithError(err).Fatal("Error establishing connection to Discord.")
		return
	}

	s.ConnDB()
	err = models.Init(s.DB)
	if err != nil {
		s.Log.WithError(err).Fatal("Error initializing database tables.")
		return
	}

	s.Log.Info("Snorlax has been woken!")
}

// Close closes the Discord session, and exits the app.
func (s *Snorlax) Close() {
	s.Log.Info("Snorlax is going to sleep.")

	err := s.Session.Close()
	if err != nil {
		s.Log.WithError(err).Error("Error closing Discord session.")
		return
	}

	err = s.DB.Close()
	if err != nil {
		s.Log.WithError(err).Error("Error closing Database connection.")
		return
	}

	s.Log.Info("Snorlax is now sleeping.")
}

// IsOwner returns whether or not a given ID is in the owners list.
func (s *Snorlax) IsOwner(id string) bool {
	for _, ownerid := range s.config.Owners {
		if ownerid == id {
			return true
		}
	}

	return false
}
