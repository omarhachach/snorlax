package moderation

import (
	"github.com/omar-h/snorlax"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Moderation"
)

// GetModule returns this module.
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "Eval contains multiple different eval commands.",
		Commands: commands,
	}
}
