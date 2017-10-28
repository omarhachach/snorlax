package music

import (
	"runtime"
	"strings"

	"github.com/dustin/go-humanize"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/utils"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Music"
)

func init() {
	playCommand := &snorlax.Command{
		Command:    ".play",
		Alias:      ".add",
		Desc:       "Play adds a YouTube video to the queue, and plays it if it's stopped.",
		Usage:      ".play <video-url>",
		ModuleName: moduleName,
		Handler:    playHandler,
	}

	stopCommand := &snorlax.Command{
		Command:    ".stop",
		Desc:       "Stops the current music stream",
		Usage:      ".stop",
		ModuleName: moduleName,
		Handler:    stopHandler,
	}

	queueCommand := &snorlax.Command{
		Command:    ".queue",
		Desc:       "Shows the music queue",
		Usage:      ".queue",
		ModuleName: moduleName,
		Handler:    queueHandler,
	}

	commands[playCommand.Command] = playCommand
	commands[stopCommand.Command] = stopCommand
	commands[queueCommand.Command] = queueCommand
}

// musicConnections maps a Guild ID to an associated voice connection.
var musicConnections = map[string]*Connection{}

var encOpts = &dca.EncodeOptions{
	Volume:           256,
	Channels:         2,
	FrameRate:        48000,
	FrameDuration:    20,
	Bitrate:          64,
	PacketLoss:       1,
	RawOutput:        true,
	Application:      dca.AudioApplicationAudio,
	CoverFormat:      "jpeg",
	CompressionLevel: 0,
	BufferedFrames:   100,
	VBR:              true,
	Threads:          0,
	AudioFilter:      "",
	Comment:          "",
}

func addToQueue(ctx *snorlax.Context, conn *Connection, song string) {
	vid, err := conn.AddYouTubeVideo(song)
	if err != nil {
		ctx.Log.WithError(err).Debug("Error adding YouTube video.")
		ctx.SendErrorMessage("YouTube link is invalid.")
		return
	}

	ctx.SendSuccessMessage("Added " + vid.Title + " to the queue.")
}

func playSong(ctx *snorlax.Context, guildID, musicChannelID, song string) {
	voice, err := ctx.Session.ChannelVoiceJoin(guildID, musicChannelID, false, true)
	if err != nil {
		ctx.Log.WithError(err).Error("Error joining voice channel.")
		return
	}
	voice.LogLevel = discordgo.LogWarning

	conn := NewConnection(voice, encOpts)
	musicConnections[guildID] = conn

	vid, err := conn.AddYouTubeVideo(song)
	if err != nil {
		ctx.Log.WithError(err).Debug("Error adding YouTube video.")
		ctx.SendErrorMessage("YouTube link is invalid.")
		return
	}

	for voice.Ready == false {
		runtime.Gosched()
	}

	ctx.SendSuccessMessage("Started playing " + vid.Title)
	err = conn.StreamMusic()
	if err != nil {
		ctx.Log.WithError(err).Error("Error starting music streaming.")
		return
	}

	conn.Close()
}

func stopHandler(ctx *snorlax.Context) {
	channel, err := ctx.State.Channel(ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting channel.")
		return
	}

	conn, ok := musicConnections[channel.GuildID]
	if !ok {
		ctx.SendErrorMessage("No music stream is playing.")
		return
	}

	conn.Close()
	ctx.SendSuccessMessage("Stopped music stream.")
}

func playHandler(ctx *snorlax.Context) {
	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	if len(parts) != 2 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	channel, err := ctx.State.Channel(ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting channel.")
		return
	}

	guild, err := ctx.State.Guild(channel.GuildID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting guild.")
		return
	}

	musicChannelID := ""
	for _, voiceState := range guild.VoiceStates {
		if voiceState.UserID == ctx.Message.Author.ID && musicChannelID == "" {
			musicChannelID = voiceState.ChannelID
		}
	}

	conn, ok := musicConnections[channel.GuildID]
	if ok {
		if musicChannelID != conn.ChannelID {
			ctx.SendErrorMessage("Please join a voice channel.")
			return
		}
		go addToQueue(ctx, conn, parts[1])
	} else {
		if musicChannelID == "" {
			ctx.SendErrorMessage("Please join a voice channel.")
			return
		}
		go playSong(ctx, channel.GuildID, musicChannelID, parts[1])
	}
}

func queueHandler(ctx *snorlax.Context) {
	channel, err := ctx.State.Channel(ctx.ChannelID)
	if err != nil {
		ctx.Log.WithError(err).Error("Error getting channel")
		return
	}

	conn, ok := musicConnections[channel.GuildID]
	if !ok {
		ctx.SendErrorMessage("No stream is playing.")
		return
	}

	queueLen := len(conn.Queue)
	if queueLen == 0 {
		ctx.SendInfoMessage("There are no queue items.")
		return
	}

	queueList := &discordgo.MessageEmbed{
		Color:  snorlax.InfoColor,
		Fields: []*discordgo.MessageEmbedField{},
	}

	n := 0
	if queueLen < 5 {
		n = queueLen
	} else {
		n = 5
	}

	for i := 0; i < n; i++ {
		queueItem := conn.Queue[i]
		queueList.Fields = append(queueList.Fields, &discordgo.MessageEmbedField{
			Name:   queueItem.Info.Title,
			Value:  queueItem.Info.Author + " - " + humanize.Time(queueItem.Info.DatePublished),
			Inline: false,
		})
	}

	ctx.SendEmbed(queueList)
}

// GetModule returns the music module.
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "The music module gives you ability to play music in a voice channel.",
		Commands: commands,
	}
}
