package snorlax

// Command holds the data and handler for a command.
type Command struct {
	Trigger        string
	Explanation    string
	ModuleName     string
	MessageHandler func()
}
