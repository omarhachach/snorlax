package rolemanager

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omarhachach/snorlax"
	"github.com/omarhachach/snorlax/utils"
)

func init() {
	setRole := &snorlax.Command{
		Command:    ".setrole",
		Alias:      ".sr",
		Desc:       "Adds a users role.",
		Usage:      ".setrole @<user> <role>",
		ModuleName: moduleName,
		Handler:    setRoleHandler,
	}

	removeRole := &snorlax.Command{
		Command:    ".removerole",
		Alias:      ".rr",
		Desc:       "Removes a users role.",
		Usage:      ".removerole @<user> <role>",
		ModuleName: moduleName,
		Handler:    removeRoleHandler,
	}

	removeAllRoles := &snorlax.Command{
		Command:    ".removeallroles",
		Alias:      ".rar",
		Desc:       "Removes all of a users roles.",
		Usage:      ".removeallroles @<user>",
		ModuleName: moduleName,
		Handler:    removeAllRolesHandler,
	}

	commands[setRole.Command] = setRole
	commands[removeRole.Command] = removeRole
	commands[removeAllRoles.Command] = removeAllRoles
}

func setRoleHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		permissions, err = ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting user permissions.")
			return
		}
	}

	if permissions&discordgo.PermissionManageRoles != 0 {
		parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
		if len(parts) != 3 {
			ctx.Log.Debugf("Wrong number of args: %#v", parts)
			return
		}

		// Get the user using the 2nd argument. (The username).
		userID := utils.ExtractUserIDFromMention(parts[1])
		user, err := ctx.Session.User(userID)
		if err != nil {
			ctx.SendErrorMessage("Username invalid.")
			ctx.Log.WithError(err).Debug("Error getting user.")
			return
		}

		channel, err := ctx.State.Channel(ctx.ChannelID)
		if err != nil {
			channel, err = ctx.Session.Channel(ctx.ChannelID)
			if err != nil {
				ctx.Log.WithError(err).Debug("Error getting guild channel.")
				return
			}
			ctx.State.ChannelAdd(channel)
		}

		roles, err := ctx.Session.GuildRoles(channel.GuildID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting guild roles.")
			return
		}

		// Check whether the role exists.
		exists := false
		roleID := ""
		for _, role := range roles {
			if !exists {
				if strings.ToLower(role.Name) == strings.ToLower(parts[2]) {
					exists = true
					roleID = role.ID
				}
			}
		}

		if !exists {
			ctx.SendErrorMessage(user.Mention() + " doesn't have role \"" + parts[2] + "\".")
			return
		}

		err = ctx.Session.GuildMemberRoleAdd(channel.GuildID, ctx.Message.Author.ID, roleID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error adding member to role.")
			return
		}

		ctx.SendSuccessMessage("Role \"" + parts[2] + "\" has been added to " + user.Mention())
	}
}

func removeRoleHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		permissions, err = ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting user permissions.")
			return
		}
	}

	// Check whether a user has the Manage Roles permission.
	if permissions&discordgo.PermissionManageRoles != 0 {
		// Get the message content and split it into arguments
		parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
		if len(parts) != 3 {
			ctx.Log.Debugf("Wrong number of args: %#v", parts)
			return
		}

		// Get the user using the 2nd argument. (The username).
		userID := utils.ExtractUserIDFromMention(parts[1])
		user, err := ctx.Session.User(userID)
		if err != nil {
			ctx.SendErrorMessage("Username invalid.")
			ctx.Log.WithError(err).Debug("Error getting user.")
			return
		}

		channel, err := ctx.State.Channel(ctx.ChannelID)
		if err != nil {
			channel, err = ctx.Session.Channel(ctx.ChannelID)
			if err != nil {
				ctx.Log.WithError(err).Debug("Error getting guild channel.")
				return
			}
			ctx.State.ChannelAdd(channel)
		}

		roles, err := ctx.Session.GuildRoles(channel.GuildID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting guild roles.")
			return
		}

		// Check whether specified role exists.
		exists := false
		roleID := ""
		for _, role := range roles {
			if !exists {
				if strings.ToLower(role.Name) == strings.ToLower(parts[2]) {
					exists = true
					roleID = role.ID
				}
			}
		}

		if !exists {
			ctx.SendErrorMessage(user.Mention() + " doesn't have role \"" + parts[2] + "\".")
			return
		}

		err = ctx.Session.GuildMemberRoleRemove(channel.GuildID, ctx.Message.Author.ID, roleID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error removing member from role.")
			return
		}

		ctx.SendSuccessMessage("Role \"" + parts[2] + "\" has been removed from " + user.Mention())
	}
}

func removeAllRolesHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		permissions, err = ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting user permissions.")
			return
		}
	}

	// Check whether a user has the Manage Roles permission.
	if permissions&discordgo.PermissionManageRoles != 0 {
		// Get the message content and split it into arguments
		parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
		if len(parts) != 2 {
			ctx.Log.Debugf("Wrong number of args: %#v", parts)
			return
		}

		// Get channel of the message (for getting GuildID)
		channel, err := ctx.State.Channel(ctx.ChannelID)
		if err != nil {
			channel, err = ctx.Session.Channel(ctx.ChannelID)
			if err != nil {
				ctx.Log.WithError(err).Debug("Error getting guild channel.")
				return
			}
			ctx.State.ChannelAdd(channel)
		}

		userID := utils.ExtractUserIDFromMention(parts[1])
		// Get Guild Member for getting roles.
		member, err := ctx.State.Member(channel.GuildID, userID)
		if err != nil {
			member, err = ctx.Session.GuildMember(channel.GuildID, userID)
			if err != nil {
				ctx.Log.WithError(err).Debug("Error getting guild member.")
				return
			}
			ctx.State.MemberAdd(member)
		}

		// Check if the user has any roles.
		userRoles := member.Roles
		if len(userRoles) <= 0 {
			ctx.SendErrorMessage(member.User.Mention() + " has no roles.")
			return
		}

		// Range over the userRoles and delete each one.
		for _, userRole := range userRoles {
			err = ctx.Session.GuildMemberRoleRemove(channel.GuildID, member.User.ID, userRole)
			if err != nil {
				ctx.Log.WithError(err).Debug("Error removing member from role.")
			}

			ctx.Log.Debug("Role deleted. ID: " + userRole)
		}

		ctx.SendSuccessMessage("All roles have been removed from " + member.User.Mention())
	}
}
