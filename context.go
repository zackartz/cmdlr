package cmdlr

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type Ctx struct {
	Session  *discordgo.Session
	Event    *discordgo.MessageCreate
	Args     *Arguments
	Router   *Router
	Command  *Command
	Database *gorm.DB
}

type ExecutionHandler func(ctx *Ctx)

func (ctx *Ctx) ResponseText(text string) error {
	_, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, text)
	return err
}

func (ctx *Ctx) RespondEmbed(embed *discordgo.MessageEmbed) error {
	_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embed)
	return err
}

func (ctx *Ctx) RespondTextEmbed(text string, embed *discordgo.MessageEmbed) error {
	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Event.ChannelID, &discordgo.MessageSend{
		Content: text,
		Embed:   embed,
	})
	return err
}
