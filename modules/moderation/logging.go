package moderation

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// This holds the different log colours.
var (
	WarnColor = 15105570
	KickColor = 15158332
	BanColor  = 12597547
)

// SendLog will send a log message.
func SendLog(s *discordgo.Session, color, points int, channel, reason, offender, caseType string) {
	s.ChannelMessageSendEmbed(channel, &discordgo.MessageEmbed{
		Title: "CASE: " + caseType,
		Color: color,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Offender",
				Value:  offender,
				Inline: false,
			},
			{
				Name:   "Issued At",
				Value:  time.Now().String(),
				Inline: false,
			},
			{
				Name:   "Points",
				Value:  strconv.Itoa(points),
				Inline: false,
			},
			{
				Name:   "Reason",
				Value:  reason,
				Inline: false,
			},
		},
	})
}

// SendWarn will send a warning log message.
func SendWarn(s *discordgo.Session, points int, channel, reason, offender string) {
	SendLog(s, WarnColor, points, channel, reason, offender, "WARN")
}

// SendKick will send a kick log message.
func SendKick(s *discordgo.Session, points int, channel, reason, offender string) {
	SendLog(s, KickColor, points, channel, reason, offender, "KICK")
}

// SendBan will send a ban log message.
func SendBan(s *discordgo.Session, points int, channel, reason, offender string) {
	SendLog(s, BanColor, points, channel, reason, offender, "BAN")
}
