package rolemanager

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/utils"
)

func init() {
	createRole := &snorlax.Command{
		Command:    ".createrole",
		Alias:      ".cr",
		Desc:       "Creates a role in the current guild.",
		ModuleName: moduleName,
		Handler:    createRoleHandler,
	}

	deleteRole := &snorlax.Command{
		Command:    ".deleterole",
		Alias:      ".dr",
		Desc:       "Deletes a role in hte current guild.",
		ModuleName: moduleName,
		Handler:    deleteRoleHandler,
	}

	commands[createRole.Command] = createRole
	commands[deleteRole.Command] = deleteRole
}

func createRoleHandler(s *snorlax.Snorlax, m *discordgo.MessageCreate) {
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
		if msgRoleName == "" || len(parts) != 3 {
			s.Log.Debug(fmt.Sprintf("Wrong number of args: %v", msgParts))
			return
		}

		channel, err := s.Session.Channel(m.ChannelID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting channel.")
			return
		}

		role, err := s.Session.GuildRoleCreate(channel.GuildID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error creating guild role.")
			return
		}

		colourIsValid, err := regexp.MatchString("^([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$", parts[1])
		if err != nil {
			s.Log.WithField("error", err).Debug("Error running regex on colour string.")
			return
		}

		if !colourIsValid {
			s.Session.ChannelMessageSend(m.ChannelID, "Colour isn't valid.")
			return
		}

		colour, err := strconv.ParseInt(parts[1], 16, 32)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error parsing colour value.")
			return
		}

		hoist, err := strconv.ParseBool(parts[2])
		if err != nil {
			s.Session.ChannelMessageSend(m.ChannelID, "Hoist value isn't a boolean (true or false).")
			s.Log.WithField("error", err).Debug("Error parsing hoist value.")
			return
		}

		role, err = s.Session.GuildRoleEdit(channel.GuildID, role.ID, msgRoleName, int(colour), hoist, 0, true)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error editing guild role.")
			return
		}

		s.Session.ChannelMessageSend(m.ChannelID, "Created role "+role.Name+"!")
	}
}

func deleteRoleHandler(s *snorlax.Snorlax, m *discordgo.MessageCreate) {
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
		if msgRoleName == "" || len(parts) != 1 {
			s.Log.Debug(fmt.Sprintf("Wrong number of args: %v", msgParts))
			return
		}

		channel, err := s.Session.Channel(m.ChannelID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting channel.")
			return
		}

		roles, err := s.Session.GuildRoles(channel.GuildID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting guild roles.")
			return
		}

		exists := false
		var roleID string
		for _, role := range roles {
			if !exists && strings.ToLower(role.Name) == strings.ToLower(msgRoleName) {
				exists = true
				roleID = role.ID
			}
		}

		if !exists {
			s.Session.ChannelMessageSend(m.ChannelID, "Role "+msgRoleName+" does not exist.")
			return
		}

		err = s.Session.GuildRoleDelete(channel.GuildID, roleID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error deleting role " + roleID + ".")
			return
		}

		s.Session.ChannelMessageSend(m.ChannelID, "Role "+msgRoleName+" has been deleted.")
	}
}
