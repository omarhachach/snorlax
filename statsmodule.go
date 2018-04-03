package snorlax

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

// This module is a conversion of the one from
// https://github.com/iopred/bruxism/blob/master/statsplugin/statsplugin.go.

func init() {
	statsModule := &Module{
		Name:     "Stats",
		Desc:     "Stats module holds a single command; `.stats`.",
		Commands: map[string]*Command{},
		Init:     statsInit,
	}

	statsCommand := &Command{
		Command:    ".stats",
		Alias:      ".info",
		Desc:       "Stats shows you statistics of the bot and runtime.",
		Usage:      ".stats",
		ModuleName: statsModule.Name,
		Handler:    statsHandler,
	}

	statsModule.Commands[statsCommand.Command] = statsCommand

	internalModules[statsModule.Name] = statsModule
}

var startTime = time.Now()

func getDuration(duration time.Duration) string {
	return fmt.Sprintf(
		"%0.2d:%0.2d:%0.2d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}

var statsMessage *discordgo.MessageEmbed
var staticFields = []*discordgo.MessageEmbedField{}

func statsInit(s *Snorlax) {
	fields := map[string]string{
		"Snorlax":   Version,
		"Go":        runtime.Version(),
		"DiscordGo": discordgo.VERSION,
		"Modules":   strconv.Itoa(len(s.Modules)),
		"Commands":  strconv.Itoa(len(s.Commands)),
	}

	for key, val := range fields {
		staticFields = append(staticFields, &discordgo.MessageEmbedField{
			Name:   key,
			Value:  val,
			Inline: true,
		})
	}

	statsMessage = &discordgo.MessageEmbed{
		Color:  InfoColor,
		Fields: staticFields,
	}

	statsReload(s)

	return
}

func statsReload(s *Snorlax) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	fields := map[string]string{
		"Uptime":      getDuration(time.Since(startTime)),
		"Memory Used": fmt.Sprintf("%s / %s", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys)),
		"Goroutines":  strconv.Itoa(runtime.NumGoroutine()),
		"Servers":     strconv.Itoa(len(s.Session.State.Guilds)),
	}

	statsMessage.Fields = staticFields
	for key, val := range fields {
		statsMessage.Fields = append(statsMessage.Fields, &discordgo.MessageEmbedField{
			Name:   key,
			Value:  val,
			Inline: true,
		})
	}
}

func statsHandler(ctx *Context) {
	statsReload(ctx.Snorlax)
	ctx.SendEmbed(statsMessage)
}
