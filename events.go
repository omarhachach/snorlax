package snorlax

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(s *Snorlax) func(sess *discordgo.Session, m *discordgo.MessageCreate) {
	return func(sess *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == sess.State.User.ID {
			return
		}

		msg := m.ContentWithMentionsReplaced()
		msgCommand := strings.Split(msg, " ")[0]

		c, ok := s.Commands[msgCommand]
		if ok {
			c.Handler(s, m)
		}
	}
}
