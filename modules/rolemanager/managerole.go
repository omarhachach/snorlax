package rolemanager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/utils"
)

func init() {
	roleHoist := &snorlax.Command{
		Command:    ".rolehoist",
		Alias:      ".rh",
		Desc:       "Role hoist changes whether or not to display a role sperately.",
		ModuleName: moduleName,
		Handler:    roleHoistHandler,
	}

	commands[roleHoist.Command] = roleHoist
}

func roleHoistHandler(s *snorlax.Snorlax, m *discordgo.MessageCreate) {
	permissions, err := s.Session.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		s.Log.WithField("error", err).Debug("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionManageRoles != 0 {
		// Get the message content and split it into arguments
		msg := m.Content
		msgParts := strings.Split(msg, " ")

		msgRoleName, parts := utils.GetStringFromParts(msgParts)
		if msgRoleName == "" || len(parts) != 2 {
			s.Log.Debug(fmt.Sprintf("Not enough arguments: %v", msgParts))
			return
		}

		channel, err := s.Session.Channel(m.ChannelID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting channel.")
			return
		}

		roles, err := s.Session.GuildRoles(channel.GuildID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting Guild Roles.")
			return
		}

		exists := false
		var role *discordgo.Role
		for _, checkRole := range roles {
			if !exists && strings.ToLower(checkRole.Name) == strings.ToLower(msgRoleName) {
				exists = true
				role = checkRole
			}
		}

		if !exists {
			s.Session.ChannelMessageSend(m.ChannelID, "Role "+msgRoleName+" does not exist!")
			return
		}

		hoist, err := strconv.ParseBool(parts[1])
		if err != nil {
			s.Session.ChannelMessageSend(m.ChannelID, "Hoist value isn't valid.")
			s.Log.WithField("error", err).Debug("Error parsing hoist boolean.")
			return
		}

		_, err = s.Session.GuildRoleEdit(channel.GuildID, role.ID, role.Name, role.Color, hoist, role.Permissions, role.Mentionable)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error editing guild role.")
			return
		}

		if hoist {
			s.Session.ChannelMessageSend(m.ChannelID, "Role "+msgRoleName+" hoisting value set to true.")
		} else {
			s.Session.ChannelMessageSend(m.ChannelID, "Role "+msgRoleName+" hoisting value set to false.")
		}
	}
}
