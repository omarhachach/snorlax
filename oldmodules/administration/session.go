package administration

import (
	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
)

func init() {
	clearStateCommand := &snorlax.Command{
		Command:    ".clearstate",
		Alias:      ".clearcache",
		Usage:      ".clearstate",
		Desc:       "Will clear the discordgo state cache.",
		ModuleName: moduleName,
		Handler:    clearStateHandler,
	}

	commands[clearStateCommand.Command] = clearStateCommand
}

func clearStateHandler(ctx *snorlax.Context) {
	if !ctx.Snorlax.IsOwner(ctx.Message.Author.ID) {
		ctx.SendErrorMessage("You have to be a bot owner to run this command.")
		return
	}

	ctx.State = discordgo.NewState()
	ctx.SendSuccessMessage("Successfully reset the discordgo state cache.")
}
