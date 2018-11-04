package utility

import (
	"github.com/omarhachach/snorlax"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Utility"
)

// GetModule returns this module.
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "Utility provides loads of utility commands, see help for more info.",
		Commands: commands,
	}
}
