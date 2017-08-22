package rolemanager

import (
	"github.com/omar-h/snorlax"
)

var (
	moduleName = "Role Manager"
	commands   = map[string]*snorlax.Command{}
)

// GetModule returns the Module
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Commands: commands,
	}
}
