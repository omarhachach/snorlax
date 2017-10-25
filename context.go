package snorlax

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Context is the context for command handlers.
// It serves as a way to unify styles and ease development.
type Context struct {
	Log           *logrus.Logger
	Session       *discordgo.Session
	MessageCreate *discordgo.MessageCreate
	Snorlax       *Snorlax
	State         *discordgo.State
	ChannelID     string
}

// This is a collection of the standard colors used for messages.
const (
	SuccessColor int = 5025616
	ErrorColor   int = 16007990
	InfoColor    int = 2201331
)

// SendEmbed sends a custom embed message.
func (ctx Context) SendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSendEmbed(ctx.ChannelID, embed)
}

// SendMessage sends an embed message with a message and color.
func (ctx Context) SendMessage(msg string, color int) (*discordgo.Message, error) {
	messageEmbed := &discordgo.MessageEmbed{
		Color:       color,
		Description: msg,
	}

	return ctx.SendEmbed(messageEmbed)
}

// SendErrorMessage is a quick way to send an error.
// -scrapped- It will suffix the message with a mention to the message creator.
func (ctx Context) SendErrorMessage(msg string) (*discordgo.Message, error) {
	//if msg[len(msg)-1:] != " " {
	//	msg += " "
	//}
	//msg += ctx.MessageCreate.Author.Mention()

	return ctx.SendMessage(msg, ErrorColor)
}

// SendSuccessMessage is a shortcut for sending a message with the success
// colors.
func (ctx Context) SendSuccessMessage(msg string) (*discordgo.Message, error) {
	return ctx.SendMessage(msg, SuccessColor)
}

// SendInfoMessage is a shortcut for sending a message with the info colors.
func (ctx Context) SendInfoMessage(msg string) (*discordgo.Message, error) {
	return ctx.SendMessage(msg, InfoColor)
}
