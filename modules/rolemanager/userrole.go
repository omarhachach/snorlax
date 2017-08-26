package rolemanager

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/utils"
)

func init() {
	setRole := snorlax.Command{
		Command:    ".setrole",
		Alias:      ".sr",
		Desc:       "Adds a users role.",
		ModuleName: moduleName,
		Handler:    setRoleHandler,
	}

	removeRole := snorlax.Command{
		Command:    ".removerole",
		Alias:      ".rr",
		Desc:       "Removes a users role.",
		ModuleName: moduleName,
		Handler:    removeRoleHandler,
	}

	removeAllRoles := snorlax.Command{
		Command:    ".removeallroles",
		Alias:      ".rar",
		Desc:       "Removes all of a users roles.",
		ModuleName: moduleName,
		Handler:    removeAllRolesHandler,
	}

	commands[setRole.Command] = &setRole
	commands[removeRole.Command] = &removeRole
	commands[removeAllRoles.Command] = &removeAllRoles
}

func setRoleHandler(s *snorlax.Snorlax, sess *discordgo.Session, m *discordgo.MessageCreate) {
	permissions, err := sess.UserChannelPermissions(m.Author.ID, m.ChannelID)
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
		userID := utils.ExtractUserIDFromMention(parts[1])
		user, err := sess.User(userID)
		if err != nil {
			sess.ChannelMessageSend(m.ChannelID, "Username invalid.")
			return
		}

		channel, err := sess.Channel(m.ChannelID)
		if err != nil {
			return
		}

		roles, err := sess.GuildRoles(channel.GuildID)
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
			sess.ChannelMessageSend(m.ChannelID, "Role \""+parts[2]+"\" does not exist.")
			return
		}
		sess.GuildMemberRoleAdd(channel.GuildID, m.Author.ID, roleID)
		sess.ChannelMessageSend(m.ChannelID, "Role \""+parts[2]+"\" has been added to "+user.Mention())
	} else {
		sess.ChannelMessageSend(m.ChannelID, "You don't have permission to do this.")
	}
}

func removeRoleHandler(s *snorlax.Snorlax, sess *discordgo.Session, m *discordgo.MessageCreate) {
	permissions, err := sess.UserChannelPermissions(m.Author.ID, m.ChannelID)
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
		userID := utils.ExtractUserIDFromMention(parts[1])
		user, err := sess.User(userID)
		if err != nil {
			sess.ChannelMessageSend(m.ChannelID, "Username invalid.")
			return
		}

		channel, err := sess.Channel(m.ChannelID)
		if err != nil {
			return
		}

		roles, err := sess.GuildRoles(channel.GuildID)
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
			sess.ChannelMessageSend(m.ChannelID, "Role \""+parts[2]+"\" does not exist.")
			return
		}
		sess.GuildMemberRoleRemove(channel.GuildID, m.Author.ID, roleID)
		sess.ChannelMessageSend(m.ChannelID, "Role \""+parts[2]+"\" has been removed from "+user.Mention())
	} else {
		sess.ChannelMessageSend(m.ChannelID, "You don't have permission to do this.")
	}
}

func removeAllRolesHandler(s *snorlax.Snorlax, sess *discordgo.Session, m *discordgo.MessageCreate) {
	permissions, err := sess.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		s.Log.Debug(fmt.Sprintf("Error getting user permissions: %v", err))
		return
	}

	if permissions&discordgo.PermissionManageRoles != 0 {
		// Get the message content and split it into arguments
		msg := m.Content
		parts := strings.Split(msg, " ")

		// Check if there are 2 arguments.
		if len(parts) != 2 {
			s.Log.Debug(fmt.Sprintf("Error running RemoveAllRoles, parts: %v", parts))
			return
		}

		// Get the user using the 2nd argument. (The username).
		userID := utils.ExtractUserIDFromMention(parts[1])
		user, err := sess.User(userID)
		if err != nil {
			sess.ChannelMessageSend(m.ChannelID, "Username invalid.")
			return
		}

		// Get channel of the message (for getting GuildID)
		channel, err := sess.Channel(m.ChannelID)
		if err != nil {
			s.Log.Debug(fmt.Sprintf("Error getting channel: %v", err))
			return
		}

		// Get Guild Member for getting roles.
		member, err := sess.GuildMember(channel.GuildID, userID)
		if err != nil {
			s.Log.Debug(fmt.Sprintf("Error getting Guild Member: %v", err))
			return
		}

		// Check if the user has any roles.
		userRoles := member.Roles
		if len(userRoles) <= 0 {
			sess.ChannelMessageSend(m.ChannelID, user.Mention()+" has no roles.")
			return
		}

		// Range over the userRoles and delete each one.
		for _, userRole := range userRoles {
			s.Log.Debug("Role deleted, ID: " + userRole)
			sess.GuildMemberRoleRemove(channel.GuildID, user.ID, userRole)
		}

		sess.ChannelMessageSend(m.ChannelID, "All roles have been removed from "+user.Mention())
	} else {
		sess.ChannelMessageSend(m.ChannelID, "You don't have permission to do this.")
	}
}
