package ping

import (
	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	log "github.com/sirupsen/logrus"
)

var (
	commands   map[string]snorlax.Command
	moduleName string
)

func init() {
	moduleName = "Ping"
	commands = make(map[string]snorlax.Command)

	pingCommand := snorlax.Command{
		Name:       "ping",
		Desc:       "Ping will respond with \"Pong!\"",
		ModuleName: moduleName,
		Handler:    ping,
	}

	commands["ping"] = pingCommand
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info(m.Content)
}

// GetModule returns the Module
func GetModule() snorlax.Module {
	return snorlax.Module{
		Name:     moduleName,
		Commands: commands,
	}
}
