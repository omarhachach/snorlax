package administration

import (
	"github.com/omar-h/snorlax"
)

func init() {
	autoDelCommand := &snorlax.Command{
		Command:    ".autodel",
		Desc:       "This will toggle the auto deletion of bot commands.",
		Usage:      ".autodel",
		ModuleName: moduleName,
		Handler:    autoDelHandler,
	}

	commands[autoDelCommand.Command] = autoDelCommand
}

func autoDelHandler(ctx *snorlax.Context) {
	if !ctx.Snorlax.IsOwner(ctx.Message.Author.ID) {
		ctx.SendErrorMessage("You have to be a bot owner to run this command.")
		return
	}

	ctx.Snorlax.Mutex.Lock()
	ctx.Snorlax.Config.AutoDelete = !ctx.Snorlax.Config.AutoDelete

	ctx.SendSuccessMessage("AutoDelete has successfully been set to %v!", ctx.Snorlax.Config.AutoDelete)
	ctx.Snorlax.Mutex.Unlock()
	err := ctx.Snorlax.Config.UpdateFile()
	if err != nil {
		ctx.Log.WithError(err).Error("Error writing to file.")
		return
	}
}
