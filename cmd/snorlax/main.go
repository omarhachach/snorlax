package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/administration"
	"github.com/omar-h/snorlax/modules/birthday"
	"github.com/omar-h/snorlax/modules/eval"
	"github.com/omar-h/snorlax/modules/moderation"
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
		administration.GetModule(),
		birthday.GetModule(),
		eval.GetModule(),
		moderation.GetModule(),
		music.GetModule(),
		ping.GetModule(),
		rolemanager.GetModule(),
	)

	bot.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	bot.Close()
}
