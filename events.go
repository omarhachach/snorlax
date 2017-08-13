package snorlax

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Content) <= 0 || m.Content[0] != '.' {
		return
	}

	msg := m.ContentWithMentionsReplaced()
	parts := strings.Split(strings.ToLower(msg), " ")
	msgCommand := strings.Replace(parts[0], ".", "", 1)

	_, ok := Commands[msgCommand]
	if !ok {
		return
	}
}
