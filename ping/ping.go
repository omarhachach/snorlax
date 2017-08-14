package ping

import (
	"github.com/omar-h/snorlax"
	"github.com/bwmarrin/discordgo"
)

var (
	commands map[string]snorlax.Command
	moduleName := "Ping"
)

func init() {
	commands := make(map[string]snorlax.Command)

	pingCommand := snorlax.Command{
		Name: "ping",
		Desc: "Ping will respond with \"Pong!\"",
		ModuleName: moduleName,
		Handler: ping
	}

	commands["ping"] = pingCommand
}

func ping(s *discordgo.Session, m *discordgo.Message) {
	
}

func GetModule() snorlax.Module {
	return snorlax.Module {
		Name: moduleName,
		Commands: commands,
	}
}
