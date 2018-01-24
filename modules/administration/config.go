package administration

import (
	"github.com/omar-h/snorlax"
	"github.com/sirupsen/logrus"
)

func init() {
	autoDelCommand := &snorlax.Command{
		Command:    ".autodel",
		Desc:       "This will toggle the auto deletion of bot commands.",
		Usage:      ".autodel",
		ModuleName: moduleName,
		Handler:    autoDelHandler,
	}

	debugCommand := &snorlax.Command{
		Command:    ".debug",
		Desc:       "This will toggle the debug mode.",
		Usage:      ".debug",
		ModuleName: moduleName,
		Handler:    debugHandler,
	}

	displayAuthor := &snorlax.Command{
		Command:    ".displayauthor",
		Desc:       "This will toggle the display author mode.",
		Usage:      ".displayauthor",
		ModuleName: moduleName,
		Handler:    displayAuthorHandler,
	}

	commands[autoDelCommand.Command] = autoDelCommand
	commands[debugCommand.Command] = debugCommand
	commands[displayAuthor.Command] = displayAuthor
}

func autoDelHandler(ctx *snorlax.Context) {
	if !ctx.Snorlax.IsOwner(ctx.Message.Author.ID) {
		ctx.SendErrorMessage("You have to be a bot owner to run this command.")
		return
	}

	ctx.Snorlax.Mutex.Lock()
	ctx.Snorlax.Config.AutoDelete = !ctx.Snorlax.Config.AutoDelete
	ctx.Snorlax.Mutex.Unlock()

	ctx.SendSuccessMessage("AutoDelete has successfully been set to %v!", ctx.Snorlax.Config.AutoDelete)

	err := ctx.Snorlax.Config.UpdateFile()
	if err != nil {
		ctx.Log.WithError(err).Error("Error writing to file.")
		return
	}
}

func debugHandler(ctx *snorlax.Context) {
	if !ctx.Snorlax.IsOwner(ctx.Message.Author.ID) {
		ctx.SendErrorMessage("You have to be a bot owner to run this command.")
		return
	}

	ctx.Snorlax.Mutex.Lock()
	ctx.Snorlax.Config.Debug = !ctx.Snorlax.Config.Debug
	if ctx.Snorlax.Config.Debug {
		ctx.Log.SetLevel(logrus.DebugLevel)
	} else {
		ctx.Log.SetLevel(logrus.ErrorLevel)
	}
	ctx.Snorlax.Mutex.Unlock()

	ctx.SendSuccessMessage("Debug mode has successfully been set to %v!", ctx.Snorlax.Config.Debug)

	err := ctx.Snorlax.Config.UpdateFile()
	if err != nil {
		ctx.Log.WithError(err).Error("Error writing to file.")
		return
	}
}

func displayAuthorHandler(ctx *snorlax.Context) {
	if !ctx.Snorlax.IsOwner(ctx.Message.Author.ID) {
		ctx.SendErrorMessage("You have to be a bot owner to run this command.")
		return
	}

	ctx.Snorlax.Mutex.Lock()
	ctx.Snorlax.Config.DisplayAuthor = !ctx.Snorlax.Config.DisplayAuthor
	ctx.Snorlax.Mutex.Unlock()

	ctx.SendSuccessMessage("Display Author mode has successfully been set to %v!", ctx.Snorlax.Config.DisplayAuthor)

	err := ctx.Snorlax.Config.UpdateFile()
	if err != nil {
		ctx.Log.WithError(err).Error("Error writing to file.")
		return
	}
}
