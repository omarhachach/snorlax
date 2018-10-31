package eval

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/robertkrimen/otto"
)

func init() {
	jsEvalCommand := &snorlax.Command{
		Command:    ".jseval",
		Desc:       "This will evaluate a JavaScript snippet and return the result.",
		Usage:      ".jseval <js-snippet>",
		ModuleName: moduleName,
		Handler:    jsEvalHandler,
	}

	commands[jsEvalCommand.Command] = jsEvalCommand
}

// vm is the JS interpreter.
var jsVM = otto.New()

func jsEvalHandler(ctx *snorlax.Context) {
	if !ctx.Snorlax.IsOwner(ctx.Message.Author.ID) {
		ctx.SendErrorMessage("You have to be a bot owner to run this command.")
		return
	}

	codeSnip := getCodeSnip(ctx.Message.Content)

	val, err := jsVM.Run(codeSnip)
	if err != nil {
		ctx.SendErrorMessage("```JS\n%s\n```", err.Error())
		return
	}

	ctx.SendEmbed(&discordgo.MessageEmbed{
		Color:       snorlax.SuccessColor,
		Description: fmt.Sprintf("```JS\n%s\n```", val),
	})
}
