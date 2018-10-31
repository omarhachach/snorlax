package moderation

import (
	"strconv"
	"strings"

	"github.com/omar-h/snorlax/modules/moderation/models"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/utils"
)

func init() {
	kickthreshold := &snorlax.Command{
		Command:    ".kickthreshold",
		Desc:       "Kick threshold will set the amount of points for a kick.",
		Usage:      ".kickthreshold <points>",
		ModuleName: moduleName,
		Handler:    kickthresholdHandler,
	}

	banthreshold := &snorlax.Command{
		Command:    ".banthreshold",
		Desc:       "Ban threshold will set the amount of points or kicks for a ban.",
		Usage:      ".banthreshold <points> <kicks>",
		ModuleName: moduleName,
		Handler:    banthresholdHandler,
	}

	commands[kickthreshold.Command] = kickthreshold
	commands[banthreshold.Command] = banthreshold
}

func kickthresholdHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionAdministrator == 0 {
		ctx.SendErrorMessage("%v doesn't have permission to ban members.", ctx.Message.Author.Mention())
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen != 2 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	points, err := strconv.Atoi(parts[1])
	if err != nil {
		ctx.Log.WithError(err).Debug("Error converting string to int.")
		ctx.SendErrorMessage("%v isn't a number.", parts[1])
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

	warnConfig.KickThreshold = points

	err = warnConfig.Insert(ctx.Snorlax.DB)
	if err != nil {
		ctx.Log.WithError(err).Error("Error inserting warn config.")
		return
	}

	ctx.SendSuccessMessage("Successfully set kickthreshold to %v!", points)
}

func banthresholdHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionAdministrator == 0 {
		ctx.SendErrorMessage("%v doesn't have permission to ban members.", ctx.Message.Author.Mention())
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen != 3 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	points, err := strconv.Atoi(parts[1])
	if err != nil {
		ctx.Log.WithError(err).Debug("Error converting string to int.")
		ctx.SendErrorMessage("%v isn't a number.", parts[1])
		return
	}

	kicks, err := strconv.Atoi(parts[2])
	if err != nil {
		ctx.Log.WithError(err).Debug("Error converting string to int.")
		ctx.SendErrorMessage("%v isn't a number.", parts[2])
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

	warnConfig.BanThreshold = points
	warnConfig.BanKickThreshold = kicks

	err = warnConfig.Insert(ctx.Snorlax.DB)
	if err != nil {
		ctx.Log.WithError(err).Error("Error inserting warn config.")
		return
	}

	ctx.SendSuccessMessage("Successfully set ban point threshold to %v and kick threshold to %v!", points, kicks)
}
