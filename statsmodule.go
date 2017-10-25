package snorlax

import (
	"bytes"
	"fmt"
	"runtime"
	"text/tabwriter"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

// This module is a conversion of the one from
// https://github.com/iopred/bruxism.

func init() {
	statsModule := &Module{
		Name:     "Stats",
		Desc:     "Stats module holds a single command; `.stats`.",
		Commands: map[string]*Command{},
		Init:     statsReload,
	}

	statsCommand := &Command{
		Command:    ".stats",
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

var statsMessage string

func statsReload(s *Snorlax) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	w.Init(buf, 0, 4, 0, ' ', 0)
	fmt.Fprint(w, "```\n")
	fmt.Fprintf(w, "Go: \t%s\n", runtime.Version())
	fmt.Fprintf(w, "Uptime: \t%s\n", getDuration(time.Since(startTime)))
	fmt.Fprintf(w, "DiscordGo: \t%s\n", discordgo.VERSION)
	fmt.Fprintf(w, "Memory Used: \t%s / %s (%s garbage collected)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc))
	fmt.Fprintf(w, "Concurrent tasks: \t%d\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "Servers: \t%d\n", len(s.Session.State.Guilds))
	fmt.Fprintf(w, "Modules: \t%d\n", len(s.Modules))
	fmt.Fprintf(w, "Commands: \t%d\n", len(s.Commands))
	fmt.Fprint(w, "```\n")

	w.Flush()
	statsMessage = buf.String()
}

func statsHandler(ctx Context) {
	statsReload(ctx.Snorlax)
	ctx.Session.ChannelMessageSend(ctx.ChannelID, statsMessage)
}
