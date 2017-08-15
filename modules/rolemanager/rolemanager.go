package rolemanager

import (
	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
)

var (
	moduleName string
	commands   map[string]snorlax.Command
)

func init() {
	moduleName = "Role Manager"
	commands = make(map[string]snorlax.Command)

	setRole := snorlax.Command{
		Name:       "setrole",
		Alias:      "sr",
		Desc:       "Sets a users role.",
		ModuleName: moduleName,
		Handler:    setRoleHandler,
	}
	commands[setRole.Name] = setRole
}

func setRoleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}
	member, err := s.GuildMember(channel.GuildID, m.Author.ID)
	if err != nil {
		return
	}
}

// GetModule returns the Module
func GetModule() snorlax.Module {
	return snorlax.Module{
		Name:     moduleName,
		Commands: commands,
	}
}
