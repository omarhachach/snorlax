package ping

import (
	"time"

	"github.com/omar-h/snorlax"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Ping"
)

func init() {
	pingCommand := &snorlax.Command{
		Command:    ".ping",
		Desc:       "Ping will respond with \"Pong!\"",
		Usage:      ".ping",
		ModuleName: moduleName,
		Handler:    ping,
	}

	commands[pingCommand.Command] = pingCommand
}

func ping(ctx snorlax.Context) {
	msgTime, err := ctx.MessageCreate.Message.Timestamp.Parse()
	if err != nil {
		ctx.Log.WithError(err).Error("ping: error parsing timestamp")
		return
	}

	msg := "Pong " + time.Since(msgTime).Round(time.Millisecond).String() + "!"
	ctx.SendMessage(msg, snorlax.InfoColor)
}

// GetModule returns the Module
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "Ping has a single command; .ping",
		Commands: commands,
	}
}
