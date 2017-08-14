package main

import (
	"flag"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/ping"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		token = flag.String("t", "", "Discord Bot Authentication Token")
	)
	flag.Parse()

	discord, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create the Discord session")
		return
	}
	bot := snorlax.NewBot(discord)

	bot.RegisterModule(ping.GetModule())
	bot.Start()
}
