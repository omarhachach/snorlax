package utility

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/omarhachach/snorlax"
	"github.com/omarhachach/snorlax/modules/moderation/models"
	"github.com/omarhachach/snorlax/utils"
)

func init() {
	portfolio := &snorlax.Command{
		Command:    ".portfolio",
		Desc:       "Portfolio will show you the portfolio URL of a given user.",
		Usage:      ".portfolio [@user]",
		ModuleName: moduleName,
		Handler:    portfolioHandler,
	}

	setPortfolio := &snorlax.Command{
		Command:    ".setportfolio",
		Desc:       "This will set the portfolio URL. Set to \"\" to reset.",
		Usage:      ".setportfolio <url> [@user]",
		ModuleName: moduleName,
		Handler:    setPortfolioHandler,
	}

	commands[portfolio.Command] = portfolio
	commands[setPortfolio.Command] = setPortfolio
}

var urlRegex = regexp.MustCompile(`(https?://)?((w{3}\.)?\w*\.\D{2,3})(\/\S*)?`)

func checkURL(url string) bool {
	return urlRegex.MatchString(url)
}

func portfolioHandler(ctx *snorlax.Context) {
	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen != 1 && partsLen != 2 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	userID := ctx.Message.Author.ID
	if partsLen == 2 {
		userID = utils.ExtractUserIDFromMention(parts[1])
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

	user, err := models.GetUser(ctx.Snorlax.DB, userID, channel.GuildID)
	if err != nil && err != models.ErrUserNotExist {
		ctx.Log.WithError(err).Error("Error getting user.")
		return
	}

	if err == models.ErrUserNotExist || user.Portfolio == "" {
		ctx.SendErrorMessage("%v doesn't have a portfolio url.", "<@"+userID+">")
		return
	}

	ctx.SendSuccessMessage("%v's portfolio is located at %v", "<@"+userID+">", user.Portfolio)
}

func setPortfolioHandler(ctx *snorlax.Context) {
	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen != 2 && partsLen != 3 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	userID := ctx.Message.Author.ID
	if partsLen == 3 {
		permissions, err := ctx.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Error("Error getting permissions.")
			return
		}

		if permissions&discordgo.PermissionAdministrator == 0 {
			ctx.SendErrorMessage("%v has to be an administrator to set someone else's portfolio url.", ctx.Message.Author.Mention())
			return
		}

		userID = utils.ExtractUserIDFromMention(parts[2])
	}

	if !checkURL(parts[1]) {
		ctx.SendErrorMessage("%v isn't a valid url.", parts[1])
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

	user, err := models.GetUser(ctx.Snorlax.DB, userID, channel.GuildID)
	if err != nil && err != models.ErrUserNotExist {
		ctx.Log.WithError(err).Error("Error getting user.")
		return
	}

	if err == models.ErrUserNotExist {
		user = &models.User{
			UserID:    userID,
			Points:    0,
			Portfolio: "",
			ServerID:  channel.GuildID,
			Kicks:     0,
		}
	}

	user.Portfolio = parts[1]

	err = user.Insert(ctx.Snorlax.DB)
	if err != nil {
		ctx.Log.WithError(err).Error("Error inserting user.")
		return
	}

	ctx.SendSuccessMessage("Successfully set %v's portfolio URL to %v.", "<@"+userID+">", parts[1])
}
