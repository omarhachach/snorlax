package birthday

import (
	"github.com/omar-h/snorlax"
)

func init() {
	setBirthdayConfigCommand := &snorlax.Command{
		Command:    ".birthdayconfig",
		Alias:      ".bdayconfig",
		Desc:       "This will set the birthday config for the server",
		Usage:      ".bdayconfig <auto-assign birthday role> <birthday role ID>",
		ModuleName: moduleName,
		Handler:    setBirthdayConfigHandler,
	}
}

func setBirthdayConfigHandler(ctx *snorlax.Context) {

}
