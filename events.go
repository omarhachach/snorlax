package snorlax

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content[0] != '.' || m.Author.ID == s.State.User.ID {
		return
	}

	msg := m.ContentWithMentionsReplaced()
	msgCommand := strings.Replace(strings.Split(strings.ToLower(msg), " ")[0], ".", "", 1)

	c, ok := Commands[msgCommand]
	if ok {
		c.Handler(s, m)
	} else {
		log.Debug("Command does not exist")
	}

}
