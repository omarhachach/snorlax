package snorlax

// Module is used to import a modular package into the bot.
// This serves to make the bot modular and expandable.
type Module struct {
	Name     string
	Commands map[string]Command
	Handlers []interface{}
}
