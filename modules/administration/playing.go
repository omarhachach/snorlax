package administration

import (
	"strings"

	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/utils"
)

func init() {
	playingCommand := &snorlax.Command{
		Command:    ".playing",
		Usage:      ".playing <message>",
		Desc:       "Will set the bot's playing message. Set to \"\" to reset.",
		ModuleName: moduleName,
		Handler:    playingHandler,
	}

	commands[playingCommand.Command] = playingCommand
}

func playingHandler(ctx *snorlax.Context) {
	if !ctx.Snorlax.IsOwner(ctx.Message.Author.ID) {
		ctx.SendErrorMessage("You have to be a bot owner to run this command.")
		return
	}

	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	if len(parts) != 2 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	err := ctx.Session.UpdateStatus(0, parts[1])
	if err != nil {
		ctx.Log.WithError(err).Error("Error updating status.")
		return
	}

	if parts[1] == "" {
		ctx.SendSuccessMessage("Successfully reset the status message.")
	} else {
		ctx.SendSuccessMessage("Successfully set the status message to %v.", parts[1])
	}
}
