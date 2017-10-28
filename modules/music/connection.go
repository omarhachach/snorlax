package music

import (
	"net/url"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

// Connection is the type for a music connection to Discord.
type Connection struct {
	GuildID         string
	EncodeOpts      *dca.EncodeOptions
	VoiceConnection *discordgo.VoiceConnection
	Queue           []*url.URL
	RWMutex         *sync.RWMutex
}

// NewConnection will return a new Connection struct.
func NewConnection(voice *discordgo.VoiceConnection, opts *dca.EncodeOptions) *Connection {
	return &Connection{
		GuildID:         voice.GuildID,
		EncodeOpts:      opts,
		VoiceConnection: voice,
		RWMutex:         &sync.RWMutex{},
	}
}

// StreamMusic will create a new encode session from the current DownloadURL
// and stream that to the VoiceConnection.
// Will block untill queue is empty.
func (c *Connection) StreamMusic() error {
	for i := 0; i < len(c.Queue); i++ {
		c.RWMutex.RLock()
		encodeSession, err := dca.EncodeFile(c.Queue[i].String(), c.EncodeOpts)
		if err != nil {
			return err
		}
		c.RWMutex.RUnlock()

		done := make(chan error)
		dca.NewStream(encodeSession, c.VoiceConnection, done)
		derr := <-done
		if derr != nil {
			return err
		}

		encodeSession.Cleanup()
	}

	return nil
}

// AddYouTubeVideo will add the download URL for a YouTube video to the queue.
func (c *Connection) AddYouTubeVideo(url string) error {
	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		return err
	}

	format := vid.Formats.Best(ytdl.FormatAudioEncodingKey)[0]
	downloadURL, err := vid.GetDownloadURL(format)
	if err != nil {
		return err
	}

	c.RWMutex.Lock()
	c.Queue = append(c.Queue, downloadURL)
	c.RWMutex.Unlock()

	return nil
}

// Close closes the VoiceConnection, stops sending speaking packet, and closes
// the EncodeSession.
func (c *Connection) Close() error {
	err := c.VoiceConnection.Speaking(false)
	if err != nil {
		return err
	}

	c.VoiceConnection.Close()
	err = c.VoiceConnection.Disconnect()
	if err != nil {
		return err
	}

	delete(musicConnections, c.GuildID)

	return nil
}
