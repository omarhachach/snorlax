package rolemanager

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
)

func init() {
	createRole := &snorlax.Command{
		Name:       "createrole",
		Alias:      "cr",
		Desc:       "Creates a role in the current guild.",
		ModuleName: moduleName,
		Handler:    createRoleHandler,
	}

	commands[createRole.Name] = createRole
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
		parts := strings.Split(msg, " ")

		if len(parts) != 4 {
			s.Log.Debug(fmt.Sprintf("Not enough arguments: %v", parts))
			return
		}

		channel, err := s.Session.Channel(m.ChannelID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting channel.")
			return
		}

		/*
			roles, err := s.Session.GuildRoles(channel.GuildID)
			if err != nil {
				s.Log.WithField("error", err).Debug("Error getting Guild Roles.")
				return
			}

			exists := false
			for _, val := range roles {
				if val.Name == parts[1] {
					exists = true
				}
			}

			if exists {
				s.Session.ChannelMessageSend(m.ChannelID, "Role "+parts[1]+" already exists!")
				return
			}
		*/

		role, err := s.Session.GuildRoleCreate(channel.GuildID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error creating GuildRole.")
			return
		}

		colourIsValid, err := regexp.MatchString("^([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$", parts[2])
		if !colourIsValid || err != nil {
			s.Session.ChannelMessageSend(m.ChannelID, "Colour isn't valid.")
			return
		}

		colour, err := strconv.ParseInt(parts[2], 16, 32)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error parsing colour value.")
			return
		}

		hoist, err := strconv.ParseBool(parts[3])
		if err != nil {
			s.Session.ChannelMessageSend(m.ChannelID, "Seperate display value isn't a boolean (true or false).")
			return
		}

		role, err = s.Session.GuildRoleEdit(channel.GuildID, role.ID, parts[1], int(colour), hoist, 0, true)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error editing guild role.")
			return
		}

		s.Session.ChannelMessageSend(m.ChannelID, "Created role "+role.Name+"!")
	}
}
