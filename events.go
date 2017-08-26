package snorlax

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(s *Snorlax) func(sess *discordgo.Session, m *discordgo.MessageCreate) {
	return func(sess *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Content[0] != '.' || m.Author.ID == s.Session.State.User.ID {
			return
		}

		msg := m.ContentWithMentionsReplaced()
		msgCommand := strings.Replace(strings.Split(strings.ToLower(msg), " ")[0], ".", "", 1)

		c, ok := s.Commands[msgCommand]
		if ok {
			c.Handler(s, m)
		} else {
			s.Log.Debug("Command " + msgCommand + " does not exist.")
		}
	}
}
