package birthday

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omarhachach/snorlax"
	"github.com/omarhachach/snorlax/modules/birthday/models"
	"github.com/omarhachach/snorlax/utils"
)

func init() {
	setBirthdayConfigCommand := &snorlax.Command{
		Command:    ".birthdayconfig",
		Alias:      ".bdayconfig",
		Desc:       "This will set the birthday config for the server",
		Usage:      ".birthdayconfig <auto-assign birthday role> <birthday role name>",
		ModuleName: moduleName,
		Handler:    setBirthdayConfigHandler,
	}

	commands[setBirthdayConfigCommand.Command] = setBirthdayConfigCommand
}

func setBirthdayConfigHandler(ctx *snorlax.Context) {
	permissions, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting user channel permissions.")
		return
	}

	if permissions&discordgo.PermissionAdministrator == 0 {
		ctx.SendErrorMessage("%v does not have permission to do this.")
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	if len(parts) != 3 {
		ctx.Log.Debug("Wrong number of parts: %v", len(parts))
		return
	}

	autoAssign, err := strconv.ParseBool(parts[1])
	if err != nil {
		ctx.SendErrorMessage("%v isn't a boolean. (True or false).")
		return
	}

	channel, err := ctx.State.Channel(ctx.ChannelID)
	if err != nil {
		channel, err = ctx.Session.Channel(ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Error("Error getting channel.")
			return
		}
		ctx.State.ChannelAdd(channel)
	}

	guildRoles, err := ctx.Session.GuildRoles(channel.GuildID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting guild roles.")
		return
	}

	exists := false
	birthdayRoleID := ""
	for i := 0; i < len(guildRoles); i++ {
		guildRole := guildRoles[i]
		if strings.ToLower(guildRole.Name) == strings.ToLower(parts[2]) {
			exists = true
			birthdayRoleID = guildRole.ID
			break
		}
	}

	if !exists {
		ctx.SendErrorMessage("Role %v does not exist.", parts[2])
	}

	bdayConfig := &models.BirthdayConfig{
		ServerID:       channel.GuildID,
		AssignRole:     autoAssign,
		BirthdayRoleID: birthdayRoleID,
	}

	err = bdayConfig.Insert(ctx.Snorlax.DB)
	if err != nil {
		ctx.Log.WithError(err).Error("Error inserting birthday config.")
		return
	}

	ctx.SendSuccessMessage("Birthday configuration has successfully been set!")
}
