package moderation

import (
	"strconv"
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

	warnConfig := &snorlax.Command{
		Command:    ".warnconfig",
		Desc:       "Will configure the warn config for the server.",
		Usage:      ".warnconfig <logwarn> <logkick> <logban>",
		ModuleName: moduleName,
		Handler:    warnConfigHandler,
	}

	commands[setWarnChannel.Command] = setWarnChannel
	commands[warnConfig.Command] = warnConfig
}

func setWarnChannelHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionAdministrator == 0 {
		ctx.SendErrorMessage("%v isn't an administrator.", ctx.Message.Author.Mention())
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
			ServerID:         channel.GuildID,
			LogChannelID:     "",
			LogWarn:          true,
			LogKick:          true,
			LogBan:           true,
			KickThreshold:    0,
			BanThreshold:     0,
			BanKickThreshold: 0,
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

func warnConfigHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionAdministrator == 0 {
		ctx.SendErrorMessage("%v isn't an administrator.", ctx.Message.Author.Mention())
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen != 4 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	logwarn, err := strconv.ParseBool(parts[1])
	if err != nil {
		ctx.SendErrorMessage("%v isn't a boolean.")
		ctx.Log.WithError(err).Debug("Error converting string to boolean.")
		return
	}

	logkick, err := strconv.ParseBool(parts[2])
	if err != nil {
		ctx.SendErrorMessage("%v isn't a boolean.")
		ctx.Log.WithError(err).Debug("Error converting string to boolean.")
		return
	}

	logban, err := strconv.ParseBool(parts[3])
	if err != nil {
		ctx.SendErrorMessage("%v isn't a boolean.")
		ctx.Log.WithError(err).Debug("Error converting string to boolean.")
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
			ServerID:         channel.GuildID,
			LogChannelID:     "",
			LogWarn:          true,
			LogKick:          true,
			LogBan:           true,
			KickThreshold:    0,
			BanThreshold:     0,
			BanKickThreshold: 0,
		}
	}

	warnConfig.LogWarn = logwarn
	warnConfig.LogKick = logkick
	warnConfig.LogBan = logban

	err = warnConfig.Insert(ctx.Snorlax.DB)
	if err != nil {
		ctx.Log.WithError(err).Error("Error inserting warn config into DB.")
		return
	}

	ctx.SendSuccessMessage("Successfully configured the warn config!")
}
