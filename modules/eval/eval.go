package eval

import (
	"strings"

	"github.com/omarhachach/snorlax"
	"github.com/robertkrimen/otto"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Eval"
)

func init() {
	vmFlushCommand := &snorlax.Command{
		Command:    ".vmflush",
		Alias:      ".vmclear",
		Desc:       "Will recreate the eval VMs.",
		Usage:      ".vmflush",
		ModuleName: moduleName,
		Handler:    vmFlushHandler,
	}

	commands[vmFlushCommand.Command] = vmFlushCommand
}

func vmFlushHandler(ctx *snorlax.Context) {
	if !ctx.Snorlax.IsOwner(ctx.Message.Author.ID) {
		ctx.SendErrorMessage("You have to be a bot owner to run this command.")
		return
	}

	jsVM = otto.New()
}

// getCodeSnip takes in the message content, and returns a single string of the
// code.
func getCodeSnip(content string) string {
	parts := strings.SplitAfter(content, " ")[1:]
	str := ""
	for i := 0; i < len(parts); i++ {
		str += parts[i]
	}
	return str
}

// GetModule returns this module.
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "Eval contains multiple different eval commands.",
		Commands: commands,
	}
}
