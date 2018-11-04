package rolemanager

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omarhachach/snorlax"
	"github.com/omarhachach/snorlax/utils"
)

func init() {
	roleHoist := &snorlax.Command{
		Command:    ".rolehoist",
		Alias:      ".rh",
		Desc:       "Role hoist changes whether or not to display a role sperately.",
		Usage:      ".rolehoist <role-name> <hoist-value>",
		ModuleName: moduleName,
		Handler:    roleHoistHandler,
	}

	roleColor := &snorlax.Command{
		Command:    ".rolecolor",
		Alias:      ".rc",
		Desc:       "Role color changes the color of a specified role.",
		Usage:      ".rolecolor <role-name> <hex-color>",
		ModuleName: moduleName,
		Handler:    roleColorHandler,
	}

	commands[roleHoist.Command] = roleHoist
	commands[roleColor.Command] = roleColor
}

func roleHoistHandler(ctx *snorlax.Context) {
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

		channel, err := ctx.State.Channel(ctx.ChannelID)
		if err != nil {
			channel, err = ctx.Session.Channel(ctx.ChannelID)
			if err != nil {
				ctx.Log.WithError(err).Debug("Error getting channel.")
				return
			}
			ctx.State.ChannelAdd(channel)
		}

		roles, err := ctx.Session.GuildRoles(channel.GuildID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting guild roles.")
			return
		}

		exists := false
		var role *discordgo.Role
		for _, checkRole := range roles {
			if !exists && strings.ToLower(checkRole.Name) == strings.ToLower(parts[1]) {
				exists = true
				role = checkRole
			}
		}

		if !exists {
			ctx.SendErrorMessage("Role " + parts[1] + " does not exist!")
			return
		}

		hoist, err := strconv.ParseBool(parts[2])
		if err != nil {
			ctx.SendErrorMessage("Hoist value isn't valid.")
			ctx.Log.WithError(err).Debug("Error parsing hoist value.")
			return
		}

		role, err = ctx.Session.GuildRoleEdit(channel.GuildID, role.ID, role.Name, role.Color, hoist, role.Permissions, role.Mentionable)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error editing guild role.")
			return
		}

		ctx.SendSuccessMessage("Hoist value for " + parts[1] + " set to " + strconv.FormatBool(hoist) + ".")
		ctx.State.RoleAdd(channel.GuildID, role)
	}
}

func roleColorHandler(ctx *snorlax.Context) {
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

		channel, err := ctx.State.Channel(ctx.ChannelID)
		if err != nil {
			channel, err = ctx.Session.Channel(ctx.ChannelID)
			if err != nil {
				ctx.Log.WithError(err).Debug("Error getting channel.")
				return
			}
			ctx.State.ChannelAdd(channel)
		}

		roles, err := ctx.Session.GuildRoles(channel.GuildID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting guild roles.")
			return
		}

		exists := false
		var role *discordgo.Role
		for _, checkRole := range roles {
			if !exists && strings.ToLower(checkRole.Name) == strings.ToLower(parts[1]) {
				exists = true
				role = checkRole
			}
		}

		if !exists {
			ctx.SendErrorMessage("Role " + parts[1] + " does not exist!")
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

		role, err = ctx.Session.GuildRoleEdit(channel.GuildID, role.ID, role.Name, int(color), role.Hoist, role.Permissions, role.Mentionable)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error editing guild role.")
			return
		}

		ctx.SendSuccessMessage("msgRoleName" + "'s colour set to " + parts[2] + ".")
		ctx.State.RoleAdd(channel.GuildID, role)
	}
}
