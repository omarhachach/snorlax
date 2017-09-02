package ping

import (
	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Ping"
)

func init() {
	pingCommand := &snorlax.Command{
		Command:    ".ping",
		Desc:       "Ping will respond with \"Pong!\"",
		Usage:      ".ping",
		ModuleName: moduleName,
		Handler:    ping,
	}

	commands[pingCommand.Command] = pingCommand
}

func ping(s *snorlax.Snorlax, m *discordgo.MessageCreate) {
	s.Session.ChannelMessageSend(m.ChannelID, "Pong! "+m.Author.Mention())
}

// GetModule returns the Module
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "Ping has a single command; .ping",
		Commands: commands,
	}
}
