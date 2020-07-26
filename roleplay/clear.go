package roleplay

import (
	"fmt"
	"strconv"
)

func ClearChannel(ctx Context) {
	channelID := ctx.TextChannel.ID
	limit := 100

	// Valid if pass number of messages
	if len(ctx.Args) > 0 {
		it, _ := strconv.Atoi(ctx.Args[0])
		limit = it
	}

	result, err := getMessages(ctx, channelID, limit)
	if err != nil {
		fmt.Println("Erro: ", err)
		ctx.Reply("Ocorreu um erro ao tentar apagar as mensagens")
		return
	}

	ctx.Reply(result)
	return
}

func getMessages(ctx Context, channelID string, limit int) (string, error) {
	var allIds []string
	var num int

	// Get messages
	messages, err := ctx.Discord.ChannelMessages(channelID, limit, "", "", "")
	if err != nil {
		return "", err
	}

	if len(messages) > 0 {
		for _, message := range messages {
			num++
			allIds = append(allIds, message.ID)
		}

		err := ctx.Discord.ChannelMessagesBulkDelete(channelID, allIds)
		if err != nil {
			return "", err
		}

		if limit == 100 {
			getMessages(ctx, channelID, limit)
		}
	}

	return fmt.Sprintf("%v mensagens apagadas", num), nil
}
