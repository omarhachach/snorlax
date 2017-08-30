package rolemanager

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/utils"
)

func init() {
	setRole := &snorlax.Command{
		Command:    ".setrole",
		Alias:      ".sr",
		Desc:       "Adds a users role.",
		ModuleName: moduleName,
		Handler:    setRoleHandler,
	}

	removeRole := &snorlax.Command{
		Command:    ".removerole",
		Alias:      ".rr",
		Desc:       "Removes a users role.",
		ModuleName: moduleName,
		Handler:    removeRoleHandler,
	}

	removeAllRoles := &snorlax.Command{
		Command:    ".removeallroles",
		Alias:      ".rar",
		Desc:       "Removes all of a users roles.",
		ModuleName: moduleName,
		Handler:    removeAllRolesHandler,
	}

	commands[setRole.Command] = setRole
	commands[removeRole.Command] = removeRole
	commands[removeAllRoles.Command] = removeAllRoles
}

func setRoleHandler(s *snorlax.Snorlax, m *discordgo.MessageCreate) {
	permissions, err := s.Session.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		s.Log.WithField("error", err).Debug("Error getting user permissions.")
		return
	}

	// Check if user has Manage Roles permission.
	if permissions&discordgo.PermissionManageRoles != 0 {
		// Get the message content and split it into arguments
		parts := utils.GetStringFromQuotes(strings.Split(m.Content, " "))
		if len(parts) != 3 {
			s.Log.Debug(fmt.Sprintf("Wrong number of args: %v", parts))
			return
		}

		// Get the user using the 2nd argument. (The username).
		userID := utils.ExtractUserIDFromMention(parts[2])
		user, err := s.Session.User(userID)
		if err != nil {
			s.Session.ChannelMessageSend(m.ChannelID, "Username invalid.")
			s.Log.WithField("error", err).Debug("Error getting user.")
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

		// Check whether the role exists.
		exists := false
		var roleID string
		for _, role := range roles {
			if !exists {
				if strings.ToLower(role.Name) == strings.ToLower(parts[1]) {
					exists = true
					roleID = role.ID
				}
			}
		}

		if !exists {
			s.Session.ChannelMessageSend(m.ChannelID, "Role \""+parts[1]+"\" does not exist.")
			return
		}
		s.Session.GuildMemberRoleAdd(channel.GuildID, m.Author.ID, roleID)
		s.Session.ChannelMessageSend(m.ChannelID, "Role \""+parts[1]+"\" has been added to "+user.Mention())
	}
}

func removeRoleHandler(s *snorlax.Snorlax, m *discordgo.MessageCreate) {
	permissions, err := s.Session.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		s.Log.WithField("error", err).Debug("Error getting user permissions.")
		return
	}

	// Check whether a user has the Manage Roles permission.
	if permissions&discordgo.PermissionManageRoles != 0 {
		// Get the message content and split it into arguments
		parts := utils.GetStringFromQuotes(strings.Split(m.Content, " "))
		if len(parts) != 3 {
			s.Log.Debug(fmt.Sprintf("Wrong number of args: %v", parts))
			return
		}

		// Get the user using the 2nd argument. (The username).
		userID := utils.ExtractUserIDFromMention(parts[2])
		user, err := s.Session.User(userID)
		if err != nil {
			s.Session.ChannelMessageSend(m.ChannelID, "Username invalid.")
			s.Log.WithField("error", err).Debug("Error getting user.")
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

		// Check whether specified role exists.
		exists := false
		var roleID string
		for _, role := range roles {
			if !exists {
				if strings.ToLower(role.Name) == strings.ToLower(parts[1]) {
					exists = true
					roleID = role.ID
				}
			}
		}

		if !exists {
			s.Session.ChannelMessageSend(m.ChannelID, "Role \""+parts[1]+"\" does not exist.")
			return
		}
		s.Session.GuildMemberRoleRemove(channel.GuildID, m.Author.ID, roleID)
		s.Session.ChannelMessageSend(m.ChannelID, "Role \""+parts[1]+"\" has been removed from "+user.Mention())
	}
}

func removeAllRolesHandler(s *snorlax.Snorlax, m *discordgo.MessageCreate) {
	permissions, err := s.Session.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		s.Log.WithField("error", err).Debug("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionManageRoles != 0 {
		// Get the message content and split it into arguments
		parts := strings.Split(m.Content, " ")

		// Check if there are 2 arguments.
		if len(parts) != 2 {
			s.Log.Debug(fmt.Sprintf("Wrong amount of args: %v", parts))
			return
		}

		// Get the user using the 2nd argument. (The username).
		userID := utils.ExtractUserIDFromMention(parts[1])
		user, err := s.Session.User(userID)
		if err != nil {
			s.Session.ChannelMessageSend(m.ChannelID, "Username invalid.")
			s.Log.WithField("error", err).Debug("Error getting user.")
			return
		}

		// Get channel of the message (for getting GuildID)
		channel, err := s.Session.Channel(m.ChannelID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting guild channel.")
			return
		}

		// Get Guild Member for getting roles.
		member, err := s.Session.GuildMember(channel.GuildID, userID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting guild member.")
			return
		}

		// Check if the user has any roles.
		userRoles := member.Roles
		if len(userRoles) <= 0 {
			s.Session.ChannelMessageSend(m.ChannelID, user.Mention()+" has no roles.")
			return
		}

		// Range over the userRoles and delete each one.
		for _, userRole := range userRoles {
			s.Log.Debug("Role deleted. ID: " + userRole)
			s.Session.GuildMemberRoleRemove(channel.GuildID, user.ID, userRole)
		}

		s.Session.ChannelMessageSend(m.ChannelID, "All roles have been removed from "+user.Mention())
	}
}
