package main

import (
	"flag"

	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/ping"
	"github.com/omar-h/snorlax/modules/rolemanager"
)

func main() {
	var (
		token = flag.String("token", "", "Discord Bot Authentication Token")
		debug = flag.Bool("debug", false, "Debug Mode")
	)
	flag.Parse()

	bot := snorlax.New(*token, &snorlax.Config{
		Debug: *debug,
	})

	bot.RegisterModule(ping.GetModule())
	bot.RegisterModule(rolemanager.GetModule())
	bot.Start()
}
