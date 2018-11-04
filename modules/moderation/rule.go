package moderation

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omarhachach/snorlax"
	"github.com/omarhachach/snorlax/modules/moderation/models"
	"github.com/omarhachach/snorlax/utils"
)

func init() {
	addRule := &snorlax.Command{
		Command:    ".addrule",
		Usage:      ".addrule <description> <points>",
		Desc:       "Will add a rule to the server's rule list.",
		ModuleName: moduleName,
		Handler:    addRuleHandler,
	}

	delRule := &snorlax.Command{
		Command:    ".delrule",
		Usage:      ".delrule <rule #>",
		Desc:       "Will remove a rule from the server's rule list.",
		ModuleName: moduleName,
		Handler:    delRuleHandler,
	}

	rules := &snorlax.Command{
		Command:    ".rules",
		Usage:      ".rules",
		Desc:       "Will display a list of rules.",
		ModuleName: moduleName,
		Handler:    rulesHandler,
	}

	commands[addRule.Command] = addRule
	commands[delRule.Command] = delRule
	commands[rules.Command] = rules
}

func addRuleHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Debug("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionAdministrator == 0 {
		ctx.SendErrorMessage("%v is not an administrator.", ctx.Message.Author.Mention())
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	if len(parts) != 3 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	points, err := strconv.Atoi(parts[2])
	if err != nil {
		ctx.Log.WithError(err).Debug("Error converting string to number.")
		ctx.SendErrorMessage("Points isn't a valid number.")
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

	rule := &models.Rule{
		Points:      points,
		ServerID:    channel.GuildID,
		Description: parts[1],
	}

	err = rule.Insert(ctx.Snorlax.DB)
	if err != nil {
		ctx.Log.WithError(err).Error("Error inserting rule.")
		return
	}

	serverRules, err := models.GetServerRules(ctx.Snorlax.DB, channel.GuildID)
	if err != nil && err != models.ErrServerRulesDontExist {
		ctx.Log.WithError(err).Error("Error getting server rules.")
		return
	}

	if err == models.ErrServerRulesDontExist {
		serverRules = &models.ServerRules{
			ServerID: channel.GuildID,
			RuleIDs:  []int{},
		}

		err = serverRules.Insert(ctx.Snorlax.DB)
		if err != nil {
			ctx.Log.WithError(err).Error("Error inserting server rules.")
			return
		}
	}

	err = serverRules.AddRule(ctx.Snorlax.DB, rule)
	if err != nil {
		ctx.Log.WithError(err).Error("Error adding rule.")
		return
	}

	ctx.SendSuccessMessage("Successfully added the rule!")
}

func delRuleHandler(ctx *snorlax.Context) {
	permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Debug("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionAdministrator == 0 {
		ctx.SendErrorMessage("%v is not an administrator.", ctx.Message.Author.Mention())
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	if len(parts) != 2 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	ruleNr, err := strconv.Atoi(parts[1])
	if err != nil {
		ctx.Log.WithError(err).Debug("Error converting string to number.")
		ctx.SendErrorMessage("%v isn't a valid rule number.", parts[1])
		return
	}

	if ruleNr <= 0 {
		ctx.SendErrorMessage("Rule # has to be greater than 0.")
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

	serverRules, err := models.GetServerRules(ctx.Snorlax.DB, channel.GuildID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting server rules.")
		return
	}

	if ruleNr > len(serverRules.RuleIDs) {
		ctx.SendErrorMessage("Rule #%v doesn't exist.", ruleNr)
		return
	}

	// Make ruleIdx be pointing to the index of the rule, rather than the actual
	// rule number.
	ruleIdx := ruleNr - 1

	err = serverRules.DelRule(ctx.Snorlax.DB, ruleIdx)
	if err != nil && err != models.ErrRuleNotExist {
		ctx.Log.WithError(err).Error("Error deleting rule.")
		return
	}

	if err == models.ErrRuleNotExist {
		ctx.SendErrorMessage("Rule %v doesn't exist.", ruleNr)
		return
	}

	ctx.SendSuccessMessage("Successfully deleted rule %v.", ruleNr)
}

var embed = discordgo.MessageEmbed{
	Color: snorlax.InfoColor,
	Fields: []*discordgo.MessageEmbedField{
		{
			Name:   "Rule # - Description - Points",
			Value:  "",
			Inline: false,
		},
	},
	Footer: &discordgo.MessageEmbedFooter{},
}

func rulesHandler(ctx *snorlax.Context) {
	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	if len(parts) != 1 {
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

	serverRules, err := models.GetServerRules(ctx.Snorlax.DB, channel.GuildID)
	if err != nil && err != models.ErrServerRulesDontExist {
		ctx.Log.WithError(err).Error("Error getting server rules.")
		return
	}

	if err == models.ErrServerRulesDontExist {
		ctx.SendInfoMessage("Server doesn't have any rules.")
		return
	}

	ruleValue := ""
	for i := 0; i < len(serverRules.RuleIDs); i++ {
		rule, err := models.GetRule(ctx.Snorlax.DB, serverRules.RuleIDs[i])
		if err != nil {
			ctx.Log.WithError(err).Error("Error getting rule.")
			return
		}

		ruleNr := strconv.Itoa(i + 1)
		rulePoints := strconv.Itoa(rule.Points)

		ruleValue = ruleValue + ruleNr + " - " + rule.Description + " - " + rulePoints + "\n"
	}

	rulesEmbed := embed
	rulesEmbed.Fields[0].Value = ruleValue
	ctx.SendEmbed(&rulesEmbed)
}
