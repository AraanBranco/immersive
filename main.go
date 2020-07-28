package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	cfg "github.com/AraanBranco/immersive/config"
	"github.com/AraanBranco/immersive/roleplay"
	"github.com/bwmarrin/discordgo"
	badgerhold "github.com/timshannon/badgerhold"
)

var (
	conf       *cfg.Config
	CmdHandler *roleplay.CommandHandler
	botId      string
	PREFIX     string
)

func init() {
	conf = cfg.LoadConfig("config.json")
	PREFIX = conf.Prefix
}

func main() {
	CmdHandler = roleplay.NewCommandHandler()
	registerCommands()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New(conf.DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session, ", err)
		return
	}

	if conf.UseSharding {
		dg.ShardID = conf.ShardID
		dg.ShardCount = conf.ShardCount
	}

	dg.AddHandler(commandHandler)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection, ", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Immersive Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botId || user.Bot {
		return
	}

	content := message.Content
	if len(content) <= len(PREFIX) {
		return
	}

	if content[:len(PREFIX)] != PREFIX {
		return
	}

	content = content[len(PREFIX):]
	if len(content) < 1 {
		return
	}

	args := strings.Fields(content)
	name := strings.ToLower(args[0])
	command, found := CmdHandler.Get(name)
	if !found {
		return
	}

	channel, err := discord.State.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel: ", err)
		return
	}

	guild, err := discord.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Error getting guild: ", err)
		return
	}

	optionsDB := badgerhold.DefaultOptions
	optionsDB.Dir = conf.LocalDB
	optionsDB.ValueDir = conf.LocalDB
	db, err := badgerhold.Open(optionsDB)
	if err != nil {
		fmt.Println("Error instance db: ", err)
	}
	defer db.Close()

	ctx := roleplay.NewContext(discord, guild, channel, user, message, conf, db, CmdHandler)
	ctx.Args = args[1:]
	c := *command
	c(*ctx)
}

func registerCommands() {
	CmdHandler.Register("ajuda", roleplay.HelpCommand, "Lista de comandos do immersive.")
	CmdHandler.Register("radio", roleplay.RadioCommand, "Gera uma nova frequência para rádio.")
	CmdHandler.Register("outfit", roleplay.OutfitCommand, "Para listar os outfits de organizações.")
	CmdHandler.Register("cidades", roleplay.GetCities, "Lista as cidades disponiveis.")
	CmdHandler.Register("chest", roleplay.ChestCommand, "Contole de itens em baú")
	CmdHandler.Register("clear", roleplay.ClearChannel, "Limpa as mensagens do canal (Admin)")
	CmdHandler.Register("contato", roleplay.ContactCommand, "Comandos criar um contato")
}
