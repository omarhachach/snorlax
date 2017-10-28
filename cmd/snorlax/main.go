package main

import (
	"flag"

	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/eval"
	"github.com/omar-h/snorlax/modules/music"
	"github.com/omar-h/snorlax/modules/ping"
	"github.com/omar-h/snorlax/modules/rolemanager"
)

var (
	token      = flag.String("token", "", "Discord Bot Authentication Token")
	debug      = flag.Bool("debug", false, "Debug Mode")
	autoDelete = flag.Bool("delete", false, "Auto Delete Mode")
)

func init() {
	flag.Parse()
}

func main() {
	bot := snorlax.New(*token, &snorlax.Config{
		Debug:     *debug,
		DeleteMsg: *autoDelete,
	})

	bot.RegisterModules(ping.GetModule(), rolemanager.GetModule(), music.GetModule(), eval.GetModule())
	bot.Start()
}
