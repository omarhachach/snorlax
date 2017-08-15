package ping

import (
	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
)

var (
	commands   map[string]*snorlax.Command
	moduleName string
)

func init() {
	moduleName = "Ping"
	commands = make(map[string]*snorlax.Command)

	pingCommand := snorlax.Command{
		Name:       "ping",
		Desc:       "Ping will respond with \"Pong!\"",
		ModuleName: moduleName,
		Handler:    ping,
	}

	commands[pingCommand.Name] = &pingCommand
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong! "+m.Author.Mention())
}

// GetModule returns the Module
func GetModule() snorlax.Module {
	return snorlax.Module{
		Name:     moduleName,
		Commands: commands,
	}
}
