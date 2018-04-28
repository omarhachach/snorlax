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
	banCommand := &snorlax.Command{
		Command:    ".ban",
		Usage:      ".ban <user> <rule #> [reason]",
		Desc:       "Will permanently remove a user from the server.",
		ModuleName: moduleName,
		Handler:    banHandler,
	}

	unbanCommand := &snorlax.Command{
		Command:    ".unban",
		Usage:      ".unban <user>",
		Desc:       "Will remove a user from the server ban list.",
		ModuleName: moduleName,
		Handler:    unbanHandler,
	}

	kickCommand := &snorlax.Command{
		Command:    ".kick",
		Usage:      ".kick <user> <rule #> [reason]",
		Desc:       "Will remove a user from the server.",
		ModuleName: moduleName,
		Handler:    kickHandler,
	}

	warnCommand := &snorlax.Command{
		Command:    ".warn",
		Usage:      ".warn <user> <rule #> [reason]",
		Desc:       "Will warn a user and assign him the points of the rule.",
		ModuleName: moduleName,
		Handler:    warnHandler,
	}

	commands[banCommand.Command] = banCommand
	commands[unbanCommand.Command] = unbanCommand
	commands[kickCommand.Command] = kickCommand
	commands[warnCommand.Command] = warnCommand
}

func banHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionBanMembers == 0 {
		ctx.SendErrorMessage("%v doesn't have permission to ban members.", ctx.Message.Author.Mention())
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen != 3 && partsLen != 4 {
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

	ruleNr, err := strconv.Atoi(parts[2])
	if err != nil {
		ctx.SendErrorMessage("%v isn't a valid rule number.", parts[2])
		ctx.Log.WithError(err).Debug("Error converting rule # to int.")
		return
	}

	if ruleNr < 0 {
		ctx.SendErrorMessage("Rule # can't be less than 0.")
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

	if ruleNr > 0 {
		serverRules, err := models.GetServerRules(ctx.Snorlax.DB, channel.GuildID)
		if err != nil && err != models.ErrServerRulesDontExist {
			ctx.Log.WithError(err).Error("Error getting server rules.")
			return
		}

		if err == models.ErrServerRulesDontExist {
			ctx.SendErrorMessage("Server doesn't have any rules.")
			return
		}

		if ruleNr > len(serverRules.RuleIDs) {
			ctx.SendErrorMessage("Server doesn't have a rule with number %v.", ruleNr)
			return
		}

		rule, err := models.GetRule(ctx.Snorlax.DB, serverRules.RuleIDs[ruleNr-1])
		if err != nil && err != models.ErrRuleNotExist {
			ctx.Log.WithError(err).Error("Error getting rule.")
			return
		}

		if err == models.ErrRuleNotExist {
			ctx.SendErrorMessage("Server doesn't have a rule with number %v.", ruleNr)
			return
		}

		// If no reason has been specified, use rule description.
		if reason == "" {
			reason = rule.Description
		}

		user, err := models.GetUser(ctx.Snorlax.DB, userID, channel.GuildID)
		if err != nil && err != models.ErrUserNotExist {
			ctx.Log.WithError(err).Error("Error getting user.")
			return
		}

		if err == models.ErrUserNotExist {
			user = &models.User{
				UserID:    userID,
				ServerID:  channel.GuildID,
				Points:    rule.Points,
				Kicks:     0,
				Portfolio: "",
			}
		} else {
			user.Points = user.Points + rule.Points
		}

		err = user.Insert(ctx.Snorlax.DB)
		if err != nil {
			ctx.Log.WithError(err).Error("Error inserting user.")
			return
		}
	}

	err = ctx.Session.GuildBanCreateWithReason(channel.GuildID, userID, reason, 0)
	if err != nil {
		ctx.SendErrorMessage("Failed to ban %v.", parts[1])
		ctx.Log.WithError(err).Debug("Failed to ban user.")
		return
	}

	ctx.SendSuccessMessage("%v has successfully been banned.", parts[1])
}

func unbanHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionBanMembers == 0 {
		ctx.SendErrorMessage("%v doesn't have permission to unban members.", ctx.Message.Author.Mention())
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	if len(parts) != 2 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	userID := utils.ExtractUserIDFromMention(parts[1])
	if userID == ctx.Message.Author.ID {
		ctx.SendErrorMessage("Can't unban yourself, obviously...")
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

func kickHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionKickMembers == 0 {
		ctx.SendErrorMessage("%v doesn't have permission to kick members.", ctx.Message.Author.Mention())
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen != 3 && partsLen != 4 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	reason := ""
	if len(parts) == 4 {
		reason = parts[3]
	}

	ruleNr, err := strconv.Atoi(parts[2])
	if err != nil {
		ctx.SendErrorMessage("%v isn't a valid rule number.", parts[2])
		ctx.Log.WithError(err).Debug("Error converting rule # to int.")
		return
	}

	if ruleNr < 0 {
		ctx.SendErrorMessage("Rule # can't be less than 0.")
		return
	}

	userID := utils.ExtractUserIDFromMention(parts[1])

	channel, err := ctx.State.Channel(ctx.ChannelID)
	if err != nil {
		channel, err = ctx.Session.Channel(ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting channel.")
			return
		}
		ctx.State.ChannelAdd(channel)
	}

	if ruleNr > 0 {
		serverRules, err := models.GetServerRules(ctx.Snorlax.DB, channel.GuildID)
		if err != nil && err != models.ErrServerRulesDontExist {
			ctx.Log.WithError(err).Error("Error getting server rules.")
			return
		}

		if err == models.ErrServerRulesDontExist {
			ctx.SendErrorMessage("Server doesn't have any rules.")
			return
		}

		if ruleNr > len(serverRules.RuleIDs) {
			ctx.SendErrorMessage("Server doesn't have a rule with number %v.", ruleNr)
			return
		}

		rule, err := models.GetRule(ctx.Snorlax.DB, serverRules.RuleIDs[ruleNr-1])
		if err != nil && err != models.ErrRuleNotExist {
			ctx.Log.WithError(err).Error("Error getting rule.")
			return
		}

		if err == models.ErrRuleNotExist {
			ctx.SendErrorMessage("Server doesn't have a rule with number %v.", ruleNr)
			return
		}

		// If no reason has been specified, use rule description.
		if reason == "" {
			reason = rule.Description
		}

		user, err := models.GetUser(ctx.Snorlax.DB, userID, channel.GuildID)
		if err != nil && err != models.ErrUserNotExist {
			ctx.Log.WithError(err).Error("Error getting user.")
			return
		}

		if err == models.ErrUserNotExist {
			user = &models.User{
				UserID:    userID,
				ServerID:  channel.GuildID,
				Points:    rule.Points,
				Kicks:     1,
				Portfolio: "",
			}
		} else {
			user.Points = user.Points + rule.Points
			user.Kicks = user.Kicks + 1
		}

		err = user.Insert(ctx.Snorlax.DB)
		if err != nil {
			ctx.Log.WithError(err).Error("Error inserting user.")
			return
		}
	}

	err = ctx.Session.GuildMemberDeleteWithReason(channel.GuildID, userID, reason)
	if err != nil {
		ctx.SendErrorMessage("Couldn't kick %v.", parts[1])
		ctx.Log.WithError(err).Debug("Couldn't kick user.")
		return
	}

	ctx.SendSuccessMessage("%v has successfully been kicked.", parts[1])
}

func warnHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionKickMembers == 0 {
		ctx.SendErrorMessage("%v doesn't have permission to kick members.", ctx.Message.Author.Mention())
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen != 3 && partsLen != 4 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	reason := ""
	if len(parts) == 4 {
		reason = parts[3]
	}

	ruleNr, err := strconv.Atoi(parts[2])
	if err != nil {
		ctx.SendErrorMessage("%v isn't a valid rule number.", parts[2])
		ctx.Log.WithError(err).Debug("Error converting rule # to int.")
		return
	}

	if ruleNr < 0 {
		ctx.SendErrorMessage("Rule # can't be less than 0.")
		return
	}

	userID := utils.ExtractUserIDFromMention(parts[1])

	channel, err := ctx.State.Channel(ctx.ChannelID)
	if err != nil {
		channel, err = ctx.Session.Channel(ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Debug("Error getting channel.")
			return
		}
		ctx.State.ChannelAdd(channel)
	}

	if ruleNr > 0 {
		serverRules, err := models.GetServerRules(ctx.Snorlax.DB, channel.GuildID)
		if err != nil && err != models.ErrServerRulesDontExist {
			ctx.Log.WithError(err).Error("Error getting server rules.")
			return
		}

		if err == models.ErrServerRulesDontExist {
			ctx.SendErrorMessage("Server doesn't have any rules.")
			return
		}

		if ruleNr > len(serverRules.RuleIDs) {
			ctx.SendErrorMessage("Server doesn't have a rule with number %v.", ruleNr)
			return
		}

		rule, err := models.GetRule(ctx.Snorlax.DB, serverRules.RuleIDs[ruleNr-1])
		if err != nil && err != models.ErrRuleNotExist {
			ctx.Log.WithError(err).Error("Error getting rule.")
			return
		}

		if err == models.ErrRuleNotExist {
			ctx.SendErrorMessage("Server doesn't have a rule with number %v.", ruleNr)
			return
		}

		// If no reason has been specified, use rule description.
		if reason == "" {
			reason = rule.Description
		}

		user, err := models.GetUser(ctx.Snorlax.DB, userID, channel.GuildID)
		if err != nil && err != models.ErrUserNotExist {
			ctx.Log.WithError(err).Error("Error getting user.")
			return
		}

		if err == models.ErrUserNotExist {
			user = &models.User{
				UserID:    userID,
				ServerID:  channel.GuildID,
				Points:    rule.Points,
				Kicks:     0,
				Portfolio: "",
			}
		} else {
			user.Points = user.Points + rule.Points
		}

		err = user.Insert(ctx.Snorlax.DB)
		if err != nil {
			ctx.Log.WithError(err).Error("Error inserting user.")
			return
		}
	}

	ctx.SendSuccessMessage("%v has successfully been warned.", parts[1])
}
