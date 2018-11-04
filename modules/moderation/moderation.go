package moderation

import (
	"github.com/omarhachach/snorlax"
	"github.com/omarhachach/snorlax/modules/moderation/models"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Moderation"
)

func moduleInit(s *snorlax.Snorlax) {
	err := models.InitRule(s.DB)
	if err != nil {
		s.Log.WithError(err).Fatal("Error initializing rule tables.")
		return
	}

	err = models.InitUsers(s.DB)
	if err != nil {
		s.Log.WithError(err).Fatal("Error initializing users tables.")
		return
	}

	err = models.InitWarnConfig(s.DB)
	if err != nil {
		s.Log.WithError(err).Fatal("Error initializing warnconfig tables.")
		return
	}
}

// GetModule returns this module.
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "Moderation contains loads of moderation tools such as ban/kick/warn.",
		Commands: commands,
		Init:     moduleInit,
	}
}
