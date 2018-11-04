package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/omarhachach/snorlax"
	"github.com/omarhachach/snorlax/modules/administration"
	"github.com/omarhachach/snorlax/modules/birthday"
	"github.com/omarhachach/snorlax/modules/eval"
	"github.com/omarhachach/snorlax/modules/moderation"
	"github.com/omarhachach/snorlax/modules/music"
	"github.com/omarhachach/snorlax/modules/ping"
	"github.com/omarhachach/snorlax/modules/rolemanager"
	"github.com/omarhachach/snorlax/modules/utility"
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
		utility.GetModule(),
	)

	bot.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	bot.Close()
}
