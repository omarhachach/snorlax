package gambling

import (
	"github.com/omar-h/snorlax"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Gambling"
)

// GetModule returns the module.
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "Gambling contains a currency system and loads of fun games.",
		Commands: commands,
	}
}
