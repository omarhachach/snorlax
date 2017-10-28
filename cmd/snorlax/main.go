package main

import (
	"flag"

	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/eval"
	"github.com/omar-h/snorlax/modules/music"
	"github.com/omar-h/snorlax/modules/ping"
	"github.com/omar-h/snorlax/modules/rolemanager"
	"github.com/sirupsen/logrus"
)

var (
	configPath = flag.String("config", "./config.json", "-config <file-path>")
)

func init() {
	flag.Parse()
}

func main() {
	config, err := snorlax.ParseConfig(*configPath)
	if err != nil {
		logrus.WithError(err).Error("Error parsing config.")
		return
	}

	bot := snorlax.New(config)

	bot.RegisterModules(
		ping.GetModule(),
		rolemanager.GetModule(),
		music.GetModule(),
		eval.GetModule(),
	)

	bot.Start()
}
