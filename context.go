package snorlax

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Context is the context for command handlers.
// It serves as a way to unify styles and ease development.
type Context struct {
	Log       *logrus.Logger
	Session   *discordgo.Session
	Message   *discordgo.MessageCreate
	Snorlax   *Snorlax
	State     *discordgo.State
	ChannelID string
}

// This is a collection of the standard colors used for messages.
const (
	SuccessColor int = 5025616
	ErrorColor   int = 16007990
	InfoColor    int = 2201331
)

// SendEmbed sends a custom embed message.
func (ctx Context) SendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	if ctx.Snorlax.Config.DisplayAuthor && embed.Author == nil {
		embed.Author = &discordgo.MessageEmbedAuthor{
			URL:     "https://www.snorlaxbot.com/",
			Name:    "Snorlax v" + Version,
			IconURL: "https://i.imgur.com/Hcoovug.png",
		}
	}

	return ctx.Session.ChannelMessageSendEmbed(ctx.ChannelID, embed)
}

// SendMessage sends an embed message with a message and color.
func (ctx Context) SendMessage(color int, title, format string, a ...interface{}) (*discordgo.Message, error) {
	messageEmbed := &discordgo.MessageEmbed{
		Color: color,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  title,
				Value: fmt.Sprintf(format, a...),
			},
		},
	}

	return ctx.SendEmbed(messageEmbed)
}

// SendErrorMessage is a quick way to send an error.
// -scrapped- It will suffix the message with a mention to the message creator.
func (ctx Context) SendErrorMessage(format string, a ...interface{}) (*discordgo.Message, error) {
	//if msg[len(msg)-1:] != " " {
	//	msg += " "
	//}
	//msg += ctx.MessageCreate.Author.Mention()

	return ctx.SendMessage(ErrorColor, "Error", format, a...)
}

// SendSuccessMessage is a shortcut for sending a message with the success
// colors.
func (ctx Context) SendSuccessMessage(format string, a ...interface{}) (*discordgo.Message, error) {
	return ctx.SendMessage(SuccessColor, "Success", format, a...)
}

// SendInfoMessage is a shortcut for sending a message with the info colors.
func (ctx Context) SendInfoMessage(format string, a ...interface{}) (*discordgo.Message, error) {
	return ctx.SendMessage(InfoColor, "Info", format, a...)
}
