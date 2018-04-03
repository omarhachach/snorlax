package snorlax

import (
	"math"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax/utils"
)

func init() {
	helpModule := &Module{
		Name:     "Help",
		Desc:     "Help gives information about a command, or displays a list of commands/modules.",
		Commands: map[string]*Command{},
		Init:     helpInit,
	}

	helpCommand := &Command{
		Command:    ".help",
		Alias:      ".h",
		Desc:       "Help shows you a help menu for a given module, or a list of modules.",
		Usage:      ".help [command/module] [page]",
		ModuleName: helpModule.Name,
		Handler:    helpHandler,
	}

	aboutCommand := &Command{
		Command:    ".about",
		Desc:       "Will tell you about the bot project.",
		Usage:      ".about",
		ModuleName: helpModule.Name,
		Handler:    aboutHandler,
	}

	helpModule.Commands[helpCommand.Command] = helpCommand
	helpModule.Commands[aboutCommand.Command] = aboutCommand

	internalModules[helpModule.Name] = helpModule
}

var (
	helpModules        *discordgo.MessageEmbed
	helpModuleCommands = map[string][]*discordgo.MessageEmbed{}
)

func helpInit(s *Snorlax) {
	botCommands = s.Commands

	helpModuleList := ""
	noDescList := ""

	for moduleName, module := range s.Modules {
		if module.Desc != "" {
			helpModuleList += moduleName + " - " + module.Desc + "\n"
		} else {
			noDescList += moduleName + "\n"
		}

		numOfPages := int(math.Ceil(float64(len(module.Commands)) / float64(20)))
		commandPages := make([]*discordgo.MessageEmbed, 0, numOfPages)

		for i := 0; i < numOfPages; i++ {
			commandPages = append(commandPages, &discordgo.MessageEmbed{
				Color: InfoColor,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Command - Description",
						Value:  "",
						Inline: false,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Use `.help [command]` for help on a command.",
				},
			})
		}

		n := 0
		for commandName, command := range module.Commands {
			n++
			page := int(math.Floor(float64(n) / float64(20)))

			commandPages[page].Fields[0].Name = "Command - Description [" + strconv.Itoa(page+1) + "/" + strconv.Itoa(numOfPages) + "]"
			commandPages[page].Fields[0].Value += commandName + " - " + command.Desc + "\n"
		}

		helpModuleCommands[strings.ToLower(moduleName)] = commandPages
	}
	helpModuleList += noDescList

	helpModules = &discordgo.MessageEmbed{
		Color: InfoColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Module Name - Description",
				Value:  helpModuleList,
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{},
	}

	return
}

var botCommands map[string]*Command

var footText = false

func helpHandler(ctx *Context) {
	parts := utils.GetStringFromQuotes(strings.Split(ctx.Message.Content, " "))

	switch len(parts) {
	case 1:
		footText = !footText
		if footText {
			helpModules.Footer.Text = "TIP: Use `.help [module] [page]` to show a list of command for a module."
		} else {
			helpModules.Footer.Text = "TIP: Surround the module's name with \"quotes\" if it spans multiple spaces."
		}
		ctx.SendEmbed(helpModules)
		return
	case 2:
		command, ok := botCommands[strings.ToLower("."+parts[1])]
		if !ok {
			commandList, ok := helpModuleCommands[strings.ToLower(parts[1])]
			if !ok {
				ctx.SendErrorMessage("Command or module \"" + parts[1] + "\" does not exist.")
				return
			}

			ctx.SendEmbed(commandList[0])
			return
		}

		ctx.SendEmbed(&discordgo.MessageEmbed{
			Color: InfoColor,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Command Usage",
					Value:  command.Usage,
					Inline: false,
				},
			},
		})
		return
	case 3:
		commandList, ok := helpModuleCommands[strings.ToLower(parts[1])]
		if !ok {
			ctx.SendErrorMessage("Module \"" + parts[1] + "\" does not exist.")
			return
		}

		pageNum, err := strconv.Atoi(parts[2])
		if err != nil {
			ctx.SendErrorMessage("\"" + parts[2] + "\"isn't a valid number.")
			ctx.Log.WithError(err).Debug("Error parsing number.")
			return
		}

		if len(commandList) < pageNum {
			ctx.SendErrorMessage("Module \"" + parts[1] + "\" only has " + strconv.Itoa(len(commandList)) + " page(s).")
			return
		}

		if pageNum < 1 {
			ctx.SendErrorMessage("Please specify a page greater than 0.")
			return
		}

		ctx.SendEmbed(commandList[pageNum-1])
		return
	default:
		ctx.Log.Debugf("Wrong number of args: %#v", parts)
		return
	}
}

var aboutEmbed = &discordgo.MessageEmbed{
	Color: InfoColor,
	Fields: []*discordgo.MessageEmbedField{
		{
			Name: "About",
			Value: "Hi, I'm Snorlax, a general purpose bot written in Go.\n\n" +
				"I am developed by Omar H., and I am open-source!\n" +
				"Support my development by contributing: https://github.com/omar-h/snorlax\n\n" +
				"Thank you very much :D",
		},
	},
}

func aboutHandler(ctx *Context) {
	ctx.SendEmbed(aboutEmbed)
}
