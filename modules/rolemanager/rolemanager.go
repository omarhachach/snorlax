package rolemanager

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
)

var (
	moduleName string
	commands   map[string]*snorlax.Command
)

func init() {
	moduleName = "Role Manager"
	commands = make(map[string]*snorlax.Command)

	setRole := snorlax.Command{
		Name:       "setrole",
		Alias:      "sr",
		Desc:       "Adds a users role.",
		ModuleName: moduleName,
		Handler:    setRoleHandler,
	}
	commands[setRole.Name] = &setRole
}

func setRoleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	permissions, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		return
	}

	if permissions&discordgo.PermissionManageRoles != 0 {
		// Get the message content and split it into arguments
		msg := m.Content
		parts := strings.Split(msg, " ")

		// Check if there are 3 arguments.
		if len(parts) != 3 {
			return
		}

		// Get the user using the 2nd argument. (The username).
		userID := strings.Replace(strings.Replace(strings.Replace(parts[1], "<", "", -1), ">", "", -1), "@", "", -1)
		user, err := s.User(userID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Username invalid.")
			return
		}

		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			return
		}

		roles, err := s.GuildRoles(channel.GuildID)
		if err != nil {
			return
		}

		exists := false
		var roleID string
		for _, role := range roles {
			if !exists {
				if role.Name == parts[2] {
					exists = true
					roleID = role.ID
				}
			}
		}

		if !exists {
			s.ChannelMessageSend(m.ChannelID, "Role \""+parts[2]+"\" does not exist.")
			return
		}
		s.GuildMemberRoleAdd(channel.GuildID, m.Author.ID, roleID)
		s.ChannelMessageSend(m.ChannelID, "Role \""+parts[2]+"\" has been added to "+user.Mention())
	} else {
		s.ChannelMessageSend(m.ChannelID, "")
	}
}

// GetModule returns the Module
func GetModule() snorlax.Module {
	return snorlax.Module{
		Name:     moduleName,
		Commands: commands,
	}
}
