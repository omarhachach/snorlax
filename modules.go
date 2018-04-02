package snorlax

// Module is used to import a modular package into the bot.
// This serves to make the bot modular and expandable.
type Module struct {
	Name     string
	Desc     string
	Commands map[string]*Command
	Init     func(*Snorlax) error
}

// Command holds the data and handler for a command.
type Command struct {
	Command    string
	Alias      string
	Desc       string
	Usage      string
	ModuleName string
	Handler    func(*Context)
}
