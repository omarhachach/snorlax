package main

import (
	"flag"

	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/ping"
	"github.com/omar-h/snorlax/modules/rolemanager"
)

func main() {
	var (
		token = flag.String("t", "", "Discord Bot Authentication Token")
	)
	flag.Parse()

	bot := snorlax.New(*token)

	bot.RegisterModule(ping.GetModule())
	bot.RegisterModule(rolemanager.GetModule())
	bot.Start()
}
