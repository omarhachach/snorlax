package rolemanager

import (
	"github.com/omarhachach/snorlax"
)

var (
	moduleName = "Role Manager"
	commands   = map[string]*snorlax.Command{}
)

// GetModule returns the Module
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "This module gives you the ability to manage roles via the bot.",
		Commands: commands,
	}
}
