package music

import (
	"io"
	"net/url"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

// QueueItem is the item used for the queue.
type QueueItem struct {
	Info *ytdl.VideoInfo
	URL  *url.URL
}

// Connection is the type for a music connection to Discord.
type Connection struct {
	GuildID         string
	ChannelID       string
	EncodeOpts      *dca.EncodeOptions
	VoiceConnection *discordgo.VoiceConnection
	Queue           []*QueueItem
	Mutex           *sync.Mutex
}

// NewConnection will return a new Connection struct.
func NewConnection(voice *discordgo.VoiceConnection, opts *dca.EncodeOptions) *Connection {
	return &Connection{
		GuildID:         voice.GuildID,
		ChannelID:       voice.ChannelID,
		EncodeOpts:      opts,
		VoiceConnection: voice,
		Mutex:           &sync.Mutex{},
	}
}

// StreamMusic will create a new encode session from the current DownloadURL
// and stream that to the VoiceConnection.
// Will block untill queue is empty.
func (c *Connection) StreamMusic() error {
	length := len(c.Queue)
	for i := 0; i < length; i++ {
		c.Mutex.Lock()
		encodeSession, err := dca.EncodeFile(c.Queue[i].URL.String(), c.EncodeOpts)
		if err != nil {
			return err
		}
		c.Mutex.Unlock()

		done := make(chan error)
		dca.NewStream(encodeSession, c.VoiceConnection, done)
		derr := <-done
		if derr != nil && derr != io.EOF {
			return derr
		}

		encodeSession.Cleanup()
		c.Mutex.Lock()
		length = len(c.Queue)
		c.Mutex.Unlock()
	}

	return nil
}

// AddYouTubeVideo will add the download URL for a YouTube video to the queue.
func (c *Connection) AddYouTubeVideo(url string) (*ytdl.VideoInfo, error) {
	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		return nil, err
	}

	format := vid.Formats.Best(ytdl.FormatAudioEncodingKey)[0]
	downloadURL, err := vid.GetDownloadURL(format)
	if err != nil {
		return nil, err
	}

	c.Mutex.Lock()
	c.Queue = append(c.Queue, &QueueItem{
		Info: vid,
		URL:  downloadURL,
	})
	c.Mutex.Unlock()

	return vid, nil
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
