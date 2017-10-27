package rolemanager

import (
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
		Usage:      ".createrole <role-name> <hex-color> <hoist-value>",
		ModuleName: moduleName,
		Handler:    createRoleHandler,
	}

	deleteRole := &snorlax.Command{
		Command:    ".deleterole",
		Alias:      ".dr",
		Desc:       "Deletes a role in hte current guild.",
		Usage:      ".deleterole <role-name>",
		ModuleName: moduleName,
		Handler:    deleteRoleHandler,
	}

	commands[createRole.Command] = createRole
	commands[deleteRole.Command] = deleteRole
}

func createRoleHandler(ctx *snorlax.Context) {
	permissions, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Debug("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionManageRoles != 0 {
		parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
		if len(parts) != 4 {
			ctx.Log.Debugf("Wrong number of args: %#v", parts)
			return
		}

		channel, err := ctx.Session.Channel(ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting channel.")
			return
		}

		role, err := ctx.Session.GuildRoleCreate(channel.GuildID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error creating guild role.")
			return
		}

		color, err := utils.HexColorToInt(parts[2])
		if err != nil {
			if err == utils.ErrColorInvalid {
				ctx.SendErrorMessage("Colour isn't valid.")
				return
			}

			ctx.Log.WithError(err).Error("Error parsing colour value.")
			return
		}

		hoist, err := strconv.ParseBool(parts[3])
		if err != nil {
			ctx.SendErrorMessage("Hoist value isn't a boolean (true or false).")
			ctx.Log.WithError(err).Debug("Error parsing hoist value.")
			return
		}

		role, err = ctx.Session.GuildRoleEdit(channel.GuildID, role.ID, parts[1], color, hoist, 0, true)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error editing guild role.")
			return
		}

		ctx.SendSuccessMessage("Created role " + parts[1] + "!")
	}
}

func deleteRoleHandler(ctx *snorlax.Context) {
	permissions, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Debug("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionManageRoles != 0 {
		parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
		if len(parts) != 4 {
			ctx.Log.Debugf("Wrong number of args: %#v", parts)
			return
		}

		channel, err := ctx.Session.Channel(ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting channel.")
			return
		}

		roles, err := ctx.Session.GuildRoles(channel.GuildID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting guild roles.")
			return
		}

		exists := false
		roleID := ""
		for _, role := range roles {
			if !exists && strings.ToLower(role.Name) == strings.ToLower(parts[1]) {
				exists = true
				roleID = role.ID
			}
		}

		if !exists {
			ctx.SendErrorMessage("Role " + parts[1] + " does not exist.")
			return
		}

		err = ctx.Session.GuildRoleDelete(channel.GuildID, roleID)
		if err != nil {
			ctx.Log.WithField("error", err).Debug("Error deleting role " + roleID + ".")
			return
		}

		ctx.SendSuccessMessage("Role " + parts[1] + " has been deleted.")
	}
}
