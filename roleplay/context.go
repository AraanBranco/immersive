package roleplay

import (
	"fmt"

	cfg "github.com/AraanBranco/immersive/config"
	"github.com/bwmarrin/discordgo"
	badgerhold "github.com/timshannon/badgerhold"
)

type Context struct {
	Discord     *discordgo.Session
	Guild       *discordgo.Guild
	TextChannel *discordgo.Channel
	User        *discordgo.User
	Message     *discordgo.MessageCreate
	Args        []string

	Conf       *cfg.Config
	DB         *badgerhold.Store
	CmdHandler *CommandHandler
}

func NewContext(discord *discordgo.Session, guild *discordgo.Guild, textChannel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate, conf *cfg.Config, db *badgerhold.Store, cmdHandler *CommandHandler) *Context {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message
	ctx.Conf = conf
	ctx.CmdHandler = cmdHandler
	ctx.DB = db
	return ctx
}

func (ctx Context) Reply(content string) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSend(ctx.TextChannel.ID, content)
	if err != nil {
		fmt.Println("Erro ao dar o reply: ", err)
		return nil
	}
	return msg
}
