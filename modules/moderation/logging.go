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

// monthMap is a map of a month's int value to its name.
var monthMap = map[int]string{
	1:  "January",
	2:  "Februar",
	3:  "March",
	4:  "April",
	5:  "May",
	6:  "June",
	7:  "July",
	8:  "August",
	9:  "September",
	10: "October",
	11: "November",
	12: "December",
}

// getTime will return a timestamp of the current time.
func getTime() string {
	now := time.Now().UTC()

	day := strconv.Itoa(now.Day())
	month := monthMap[int(now.Month())]
	year := strconv.Itoa(now.Year())

	return day + " " + month + " " + year
}

// SendLog will send a log message.
func SendLog(s *discordgo.Session, color, points int, channel, reason, offender, caseType, imageURL string) {
	s.ChannelMessageSendEmbed(channel, &discordgo.MessageEmbed{
		Title: "CASE: " + caseType,
		Color: color,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Offender",
				Value:  offender,
				Inline: true,
			},
			{
				Name:   "Points",
				Value:  strconv.Itoa(points),
				Inline: true,
			},
			{
				Name:   "Reason",
				Value:  reason,
				Inline: true,
			},
			{
				Name:   "Issued At",
				Value:  getTime(),
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    imageURL,
			Width:  128,
			Height: 128,
		},
	})
}

// SendWarn will send a warning log message.
func SendWarn(s *discordgo.Session, points int, channel, reason, offender, imageURL string) {
	SendLog(s, WarnColor, points, channel, reason, offender, "WARN", imageURL)
}

// SendKick will send a kick log message.
func SendKick(s *discordgo.Session, points int, channel, reason, offender, imageURL string) {
	SendLog(s, KickColor, points, channel, reason, offender, "KICK", imageURL)
}

// SendBan will send a ban log message.
func SendBan(s *discordgo.Session, points int, channel, reason, offender, imageURL string) {
	SendLog(s, BanColor, points, channel, reason, offender, "BAN", imageURL)
}
