package moderation

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/moderation/models"
	"github.com/omar-h/snorlax/utils"
)

func init() {
	setWarnChannel := &snorlax.Command{
		Command:    ".setlogchannel",
		Alias:      ".setlogchnl",
		Desc:       "Will set the warn/kick/ban logging channel to the current one.",
		Usage:      ".setlogchannel",
		ModuleName: moduleName,
		Handler:    setWarnChannelHandler,
	}

	commands[setWarnChannel.Command] = setWarnChannel
}

func setWarnChannelHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionAdministrator != 0 {
		ctx.SendErrorMessage("%v doesn't have permission to ban members.", ctx.Message.Author.Mention())
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen != 1 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
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

	warnConfig, err := models.GetWarnConfig(ctx.Snorlax.DB, channel.GuildID)
	if err != nil && err != models.ErrWarnConfigNotExist {
		ctx.Log.WithError(err).Error("Error getting warn config.")
		return
	}

	if err == models.ErrWarnConfigNotExist {
		warnConfig = &models.WarnConfig{
			ServerID:     channel.GuildID,
			LogChannelID: "",
			LogWarn:      true,
			LogKick:      true,
			LogBan:       true,
		}
	}

	warnConfig.LogChannelID = ctx.ChannelID

	err = warnConfig.Insert(ctx.Snorlax.DB)
	if err != nil {
		ctx.Log.WithError(err).Error("Error inserting warn config into DB.")
		return
	}

	ctx.SendSuccessMessage("Successfully set warn channel ID to %v!", ctx.ChannelID)
}
