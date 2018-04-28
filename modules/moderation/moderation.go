package moderation

import (
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/moderation/models"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Moderation"
)

func moduleInit(s *snorlax.Snorlax) {
	err := models.InitRule(s.DB)
	if err != nil {
		s.Log.WithError(err).Fatal("Error initializing moderation module.")
		return
	}
}

// GetModule returns this module.
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "Eval contains multiple different eval commands.",
		Commands: commands,
		Init:     moduleInit,
	}
}
