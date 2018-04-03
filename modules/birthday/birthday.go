package birthday

import (
	"strings"
	"time"

	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/birthday/models"
	"github.com/omar-h/snorlax/utils"
)

var (
	commands   = map[string]*snorlax.Command{}
	moduleName = "Birthday"
)

func init() {
	setBirthdayCommand := &snorlax.Command{
		Command:    ".setbirthday",
		Alias:      ".setbday",
		Desc:       "Will set a birthday.",
		Usage:      ".setbday <MM/DD> [user]",
		ModuleName: moduleName,
		Handler:    setBirthdayHandler,
	}

	birthdayCommand := &snorlax.Command{
		Command:    ".birthday",
		Alias:      ".bday",
		Desc:       "Will display a birthday.",
		Usage:      ".bday [user]",
		ModuleName: moduleName,
		Handler:    birthdayHandler,
	}

	commands[setBirthdayCommand.Command] = setBirthdayCommand
	commands[birthdayCommand.Command] = birthdayCommand
}

func birthdayInit(s *snorlax.Snorlax) {
	err := models.BirthdayInit(s.DB)
	if err != nil {
		s.Log.WithError(err).Error("Error initializing birthday table.")
		return
	}

	err = models.BirthdayConfigInit(s.DB)
	if err != nil {
		s.Log.WithError(err).Error("Error initalizing birthday config init.")
	}

	birthdayTimer(s)
	return
}

// cachedDate provides the date which it was last time the birthdayTimer
// function ran. This will make sure we don't run until it changes.
var cachedDate = time.Now().Day()

func birthdayTimer(s *snorlax.Snorlax) {
	day := time.Now().Day()

	if day == cachedDate {
		time.Sleep(1 * time.Hour)
		birthdayTimer(s)
		return
	}
	cachedDate = day

	giveBirthdayRoles(s)
	removeBirthdayRoles(s)

	time.Sleep(1 * time.Hour)
	birthdayTimer(s)
}

func setBirthdayHandler(ctx *snorlax.Context) {
	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen < 2 || partsLen > 3 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	if partsLen == 2 {
		parts = append(parts, ctx.Message.Author.ID)
	} else {
		parts[2] = utils.ExtractUserIDFromMention(parts[2])
	}

	channel, err := ctx.State.Channel(ctx.ChannelID)
	if err != nil {
		channel, err = ctx.Session.Channel(ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Error("Failed to get channel.")
		}
		ctx.State.ChannelAdd(channel)
	}

	bday := &models.Birthday{
		UserID:   parts[2],
		Birthday: parts[1],
		ServerID: channel.GuildID,
	}

	err = bday.Insert(ctx.Snorlax.DB)
	if err != nil {
		ctx.Log.WithError(err).Error("Error inserting birthday.")
		ctx.SendErrorMessage("Error setting birthday.")
		return
	}

	ctx.SendSuccessMessage("Successfully set your birthday to %v!", parts[1])
}

func birthdayHandler(ctx *snorlax.Context) {
	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))
	partsLen := len(parts)
	if partsLen < 1 || partsLen > 2 {
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}

	if partsLen == 1 {
		parts = append(parts, ctx.Message.Author.ID)
	} else {
		parts[2] = utils.ExtractUserIDFromMention(parts[2])
	}

	channel, err := ctx.State.Channel(ctx.ChannelID)
	if err != nil {
		channel, err = ctx.Session.Channel(ctx.ChannelID)
		if err != nil {
			ctx.Log.WithError(err).Error("Error getting channel.")
			return
		}
		ctx.State.ChannelAdd(channel)
	}

	birthday, err := models.GetBirthday(ctx.Snorlax.DB, parts[0], channel.GuildID)
	if err != nil && err != models.ErrNoBirthdayFound {
		ctx.Log.WithError(err).Error("Error getting birthday.")
		return
	}

	if err == models.ErrNoBirthdayFound {
		ctx.SendErrorMessage("No birthday was found.")
		return
	}

	ctx.SendSuccessMessage("<@%v>'s birthday is %v.", birthday.UserID, birthday.Birthday)
}

// GetModule returns the Module
func GetModule() *snorlax.Module {
	return &snorlax.Module{
		Name:     moduleName,
		Desc:     "The birthday allows for handling and assigning special birthday things.",
		Commands: commands,
		Init:     birthdayInit,
	}
}
