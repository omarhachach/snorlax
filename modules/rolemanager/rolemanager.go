package rolemanager

import (
	"github.com/omar-h/snorlax"
)

var (
	moduleName string
	commands   map[string]*snorlax.Command
)

func init() {
	moduleName = "Role Manager"
	commands = make(map[string]*snorlax.Command)
}

// GetModule returns the Module
func GetModule() snorlax.Module {
	return snorlax.Module{
		Name:     moduleName,
		Commands: commands,
	}
}
