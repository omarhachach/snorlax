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
			if s.config.DeleteMsg {
				go func() {
					err := s.Session.ChannelMessageDelete(m.ChannelID, m.ID)
					if err != nil {
						s.Log.WithError(err).Error("Error automatically deleting message.")
					}
				}()
			}

			go c.Handler(Context{
				Log:       s.Log,
				Session:   sess,
				Snorlax:   s,
				State:     sess.State,
				Message:   m,
				ChannelID: m.ChannelID,
			})
		}
	}
}
