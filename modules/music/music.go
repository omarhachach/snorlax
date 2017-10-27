package music

import (
	"runtime"
	"strings"

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
		Desc:       "Play plays a YouTube video in a voice channel.",
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

	commands[playCommand.Command] = playCommand
	commands[stopCommand.Command] = stopCommand
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

	if musicChannelID == "" {
		ctx.SendErrorMessage("Please join a voice channel.")
		return
	}

	voice, err := ctx.Session.ChannelVoiceJoin(channel.GuildID, musicChannelID, false, true)
	if err != nil {
		ctx.Log.WithError(err).Error("Error joining voice channel.")
		return
	}
	voice.LogLevel = discordgo.LogWarning

	conn := NewConnection(voice, encOpts)
	musicConnections[channel.GuildID] = conn

	err = conn.AddYouTubeVideo(parts[1])
	if err != nil {
		ctx.Log.WithError(err).Debug("Error adding YouTube video.")
		ctx.SendErrorMessage("YouTube link is invalid.")
		return
	}

	for voice.Ready == false {
		runtime.Gosched()
	}

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

// GetModule returns the music module.
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "The music module gives you ability to play music in a voice channel.",
		Commands: commands,
	}
}
