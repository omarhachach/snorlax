package moderation

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/utils"
)

func init() {
	banCommand := &snorlax.Command{
		Command:    ".ban",
		Usage:      ".ban <user> [reason]",
		Desc:       "Will permanently remove a user from the server.",
		ModuleName: moduleName,
		Handler:    banHandler,
	}

	// tempBanCommand := &snorlax.Command{
	// 	Command:    ".tempban",
	// 	Usage:      ".tempban <user> <reason>",
	// 	Alias:      ".tempkick",
	// 	Desc:       "Will temporarily remove a user from the server. 0 = perm.",
	// 	ModuleName: moduleName,
	// 	Handler:    tempBanHandler,
	// }

	unbanCommand := &snorlax.Command{
		Command:    ".unban",
		Usage:      ".unban <user>",
		Desc:       "Will remove a user from the server ban list.",
		ModuleName: moduleName,
		Handler:    unbanHandler,
	}

	kickCommand := &snorlax.Command{
		Command:    ".kick",
		Usage:      ".kick <user> <reason>",
		Desc:       "Will remove a user from the server.",
		ModuleName: moduleName,
		Handler:    kickHandler,
	}

	commands[banCommand.Command] = banCommand
	// commands[tempBanCommand.Command] = tempBanCommand
	commands[unbanCommand.Command] = unbanCommand
	commands[kickCommand.Command] = kickCommand
}

func banUser(ctx *snorlax.Context, userID, reason string, time int) (bool, error) {
	channel, err := ctx.Session.Channel(ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Debug("Error getting channel.")
		return false, nil // Don't want to return error, as error is handled.
	}

	return true, ctx.Session.GuildBanCreateWithReason(channel.GuildID, userID, reason, time)
}

func banHandler(ctx *snorlax.Context) {
	permissions, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Debug("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionBanMembers != 0 {
		parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
		partsLen := len(parts)
		if partsLen != 2 && partsLen != 3 {
			ctx.Log.Debugf("Wrong number of args: %#v", parts)
			return
		}

		userID := utils.ExtractUserIDFromMention(parts[1])
		if userID == ctx.Message.Author.ID {
			ctx.SendErrorMessage("Can't ban yourself.")
			return
		}

		reason := ""
		if len(parts) == 3 {
			reason = parts[2]
		}
		ok, err := banUser(ctx, userID, reason, 0)
		if ok == false {
			return
		}

		if err != nil {
			ctx.SendErrorMessage("Failed to ban %v.", parts[1])
			ctx.Log.WithError(err).Debug("Failed to ban user.")
			return
		}
		ctx.SendSuccessMessage("%v has successfully been banned.", parts[1])
	}
}

// func tempBanHandler(ctx *snorlax.Context) {
// 	permissions, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
// 	if err != nil {
// 		ctx.Log.WithError(err).Debug("Error getting user permissions.")
// 		return
// 	}

// 	if permissions&discordgo.PermissionBanMembers != 0 {
// 		parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
// 		if len(parts) != 4 {
// 			ctx.Log.Debugf("Wrong number of args: %#v", parts)
// 			return
// 		}

// 		time, err := strconv.Atoi(parts[3])
// 		if err != nil {
// 			ctx.SendErrorMessage("%v is not a number.", parts[3])
// 			ctx.Log.WithError(err).Debug("Failed to convert string to number.")
// 			return
// 		}

// 		if time < 0 {
// 			ctx.SendErrorMessage("Time cannot be less than 0.")
// 			return
// 		}

// 		ok, err := banUser(ctx, utils.ExtractUserIDFromMention(parts[1]), parts[2], time)
// 		if ok == false {
// 			return
// 		}

// 		if err != nil {
// 			ctx.SendErrorMessage("Failed to tempban %v.", parts[1])
// 			ctx.Log.WithError(err).Debug("Failed to tempban user.")
// 			return
// 		}
// 	}
// }

func unbanHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		permissions, err = ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting user permissions.")
			return
		}
	}

	if permissions&discordgo.PermissionBanMembers != 0 {
		parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
		if len(parts) != 2 {
			ctx.Log.Debugf("Wrong number of args: %#v", parts)
			return
		}

		userID := utils.ExtractUserIDFromMention(parts[1])
		if userID == ctx.Message.Author.ID {
			ctx.SendErrorMessage("Can't unban yourself. (How is this even happening?).")
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

		bans, err := ctx.Session.GuildBans(channel.GuildID)
		if err != nil {
			ctx.Log.WithError(err).Error("Error getting guild bans.")
			return
		}

		exists := false
		for _, ban := range bans {
			if ban.User.ID == userID {
				exists = true
			}
		}

		if !exists {
			ctx.SendErrorMessage("%v isn't banned.", parts[1])
			return
		}

		err = ctx.Session.GuildBanDelete(channel.GuildID, userID)
		if err != nil {
			ctx.SendErrorMessage("Failed to unban %v.", parts[1])
			ctx.Log.WithError(err).Debug("Failed to unban user.")
			return
		}
		ctx.SendSuccessMessage("%v has successfully been unbanned.", parts[1])
	}
}

func kickHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		permissions, err = ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting user permissions.")
			return
		}
	}

	if permissions&discordgo.PermissionBanMembers != 0 {
		parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
		partsLen := len(parts)
		if partsLen != 2 && partsLen != 3 {
			ctx.Log.Debugf("Wrong number of args: %#v", parts)
			return
		}

		reason := ""
		if len(parts) == 3 {
			reason = parts[2]
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

		err = ctx.Session.GuildMemberDeleteWithReason(channel.GuildID, utils.ExtractUserIDFromMention(parts[1]), reason)
		if err != nil {
			ctx.SendErrorMessage("Couldn't kick %v.", parts[1])
			ctx.Log.WithError(err).Debug("Couldn't kick user.")
			return
		}
		ctx.SendSuccessMessage("%v has successfully been kicked.", parts[1])
	}

}
