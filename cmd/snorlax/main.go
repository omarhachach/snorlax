package main

import (
	"flag"

	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/ping"
	"github.com/omar-h/snorlax/modules/rolemanager"
)

var (
	token = flag.String("token", "", "Discord Bot Authentication Token")
	debug = flag.Bool("debug", false, "Debug Mode")
)

func init() {
	flag.Parse()
}

func main() {
	bot := snorlax.New(*token, &snorlax.Config{
		Debug: *debug,
	})

	bot.RegisterModule(ping.GetModule())
	bot.RegisterModule(rolemanager.GetModule())
	bot.Start()
}
