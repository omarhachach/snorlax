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
	addRule := &snorlax.Command{
		Command:    ".addrule",
		Usage:      ".addrule <rule> <points>",
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

	points, err := strconv.Atoi(parts[1])
	if err != nil {
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
		Description: parts[2],
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
	}

	err = serverRules.AddRule(ctx.Snorlax.DB, rule)
	if err != nil {
		ctx.Log.WithError(err).Error("Error adding rule.")
		return
	}

	ctx.SendSuccessMessage("Successfully added the rule!")
}

func delRuleHandler(ctx *snorlax.Context) {

}

func rulesHandler(ctx *snorlax.Context) {

}
