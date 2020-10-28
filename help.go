package cmdlr

import (
	"github.com/bwmarrin/discordgo"
	"math"
	"strconv"
	"strings"
	"time"
)

func (r *Router) RegisterDefaultHelpCommand(session *discordgo.Session) {
	r.InitializeStorage("hdl_helpMessages")

	session.AddHandler(func(session *discordgo.Session, event *discordgo.MessageReactionAdd) {
		channelID := event.ChannelID
		messageID := event.MessageID
		userID := event.UserID

		if event.UserID == session.State.User.ID {
			return
		}

		rawPage, ok := r.Storage["hdl_helpMessages"].Get(channelID + ":" + messageID + ":" + userID)
		if !ok {
			return
		}

		page := rawPage.(int)
		if page <= 0 {
			return
		}

		reactionName := event.Emoji.Name
		switch reactionName {
		case "⬅️":
			embed, newPage := renderDefaultGeneralHelpEmbed(r, page-1)
			page = newPage
			session.ChannelMessageEditEmbed(channelID, messageID, embed)

			session.MessageReactionRemove(channelID, messageID, reactionName, userID)
			break
		case "❌":
			session.ChannelMessageDelete(channelID, messageID)
			break
		case "➡️":
			embed, newPage := renderDefaultGeneralHelpEmbed(r, page+1)
			page = newPage
			session.ChannelMessageEditEmbed(channelID, messageID, embed)

			session.MessageReactionRemove(channelID, messageID, reactionName, userID)
			break
		}

		r.Storage["hdl_helpMessages"].Set(channelID+":"+messageID+":"+userID, page)
	})

	r.RegisterCmd(&Command{
		Name:        "help",
		Description: "Lists all the available commands or displays some information about a specific command",
		Usage:       "help [command name]",
		Example:     "help yourCommand",
		IgnoreCase:  true,
		Handler:     generalHelpCommand,
	})

}

func generalHelpCommand(ctx *Ctx) {
	if ctx.Args.Amount() > 0 {
		specificHelpCommand(ctx)
		return
	}

	channelID := ctx.Event.ChannelID
	session := ctx.Session

	embed, _ := renderDefaultGeneralHelpEmbed(ctx.Router, 1)
	message, _ := ctx.Session.ChannelMessageSendEmbed(channelID, embed)

	session.MessageReactionAdd(channelID, message.ID, "⬅️")
	session.MessageReactionAdd(channelID, message.ID, "❌")
	session.MessageReactionAdd(channelID, message.ID, "➡️")

	ctx.Router.Storage["hdl_helpMessages"].Set(channelID+":"+message.ID+":"+ctx.Event.Author.ID, 1)
}

func specificHelpCommand(ctx *Ctx) {
	// Define the command names
	commandNames := strings.Split(ctx.Args.Raw(), " ")

	// Define the command
	var command *Command
	for index, commandName := range commandNames {
		if index == 0 {
			command = ctx.Router.GetCmd(commandName)
			continue
		}
		command = command.GetSubCommand(commandName)
	}

	// Send the help embed
	ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, renderDefaultSpecificHelpEmbed(ctx, command))
}

func renderDefaultGeneralHelpEmbed(r *Router, page int) (*discordgo.MessageEmbed, int) {
	commands := r.Commands
	prefix := r.Prefixes[0]

	pageAmount := int(math.Ceil(float64(len(commands)) / 5))
	if page > pageAmount {
		page = pageAmount
	}
	if page <= 0 {
		page = 1
	}

	startingIndex := (page - 1) * 5
	endingIndex := startingIndex + 5
	if page == pageAmount {
		endingIndex = len(commands)
	}
	displayCommands := commands[startingIndex:endingIndex]

	fields := make([]*discordgo.MessageEmbedField, len(displayCommands))
	for index, command := range displayCommands {
		fields[index] = &discordgo.MessageEmbedField{
			Name:   command.Name,
			Value:  "`" + command.Description + "`",
			Inline: false,
		}
	}

	return &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "Command List (Page " + strconv.Itoa(page) + "/" + strconv.Itoa(pageAmount) + ")",
		Description: "These are all the available commands. Type `" + prefix + "help <command_name> to find out more about a specific command.",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0xffff00,
		Fields:      fields,
	}, page
}

func renderDefaultSpecificHelpEmbed(ctx *Ctx, command *Command) *discordgo.MessageEmbed {
	prefix := ctx.Router.Prefixes[0]

	if command == nil {
		return &discordgo.MessageEmbed{
			Type:      "rich",
			Title:     "Error",
			Timestamp: time.Now().Format(time.RFC3339),
			Color:     0xff0000,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Message",
					Value:  "```The given command doesn't exist. Type `" + prefix + "help` for a list of available commands.```",
					Inline: false,
				},
			},
		}
	}

	subCommands := "No sub commands"
	if len(command.SubCommands) > 0 {
		subCommandNames := make([]string, len(command.SubCommands))
		for index, subCommand := range command.SubCommands {
			subCommandNames[index] = subCommand.Name
		}
		subCommands = "`" + strings.Join(subCommandNames, "`, `") + "`"
	}

	aliases := "No aliases"
	if len(command.Aliases) > 0 {
		aliases = "`" + strings.Join(command.Aliases, "`, `") + "`"
	}

	return &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "Command Information",
		Description: "Displaying the information for the `" + command.Name + "` command.",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0xffff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Name",
				Value:  "`" + command.Name + "`",
				Inline: false,
			},
			{
				Name:   "Sub Commands",
				Value:  subCommands,
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  aliases,
				Inline: false,
			},
			{
				Name:   "Description",
				Value:  "```" + command.Description + "```",
				Inline: false,
			},
			{
				Name:   "Usage",
				Value:  "```" + prefix + command.Usage + "```",
				Inline: false,
			},
			{
				Name:   "Example",
				Value:  "```" + prefix + command.Example + "```",
				Inline: false,
			},
		},
	}
}
