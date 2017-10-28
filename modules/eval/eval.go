package eval

import (
	"strings"

	"github.com/omar-h/snorlax"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Eval"
)

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
