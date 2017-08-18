package main

import (
	"flag"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/ping"
	"github.com/omar-h/snorlax/modules/rolemanager"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		token = flag.String("t", "", "Discord Bot Authentication Token")
		debug = flag.Bool("debug", true, "Debug Mode")
	)
	flag.Parse()

	discord, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create the Discord session")
		return
	}
	bot := snorlax.NewBot(discord, &snorlax.Config{
		Debug: *debug,
	})

	bot.RegisterModule(ping.GetModule())
	bot.RegisterModule(rolemanager.GetModule())
	bot.Start()
}
