package roleplay

import "fmt"

func LogChestToOwner(ctx Context, message string) {
	dm, err := ctx.Discord.UserChannelCreate(ctx.Config.GetEnvConfString("OwnerId"))
	if err != nil {
		fmt.Println("Erro ao enviar DM: ", err.Error())
	}

	msg := fmt.Sprintf("Mensagem enviada por: %v \n %v", ctx.Message.Author.Mention(), message)

	// Send message private
	_, err = ctx.Discord.ChannelMessageSend(dm.ID, msg)
	if err != nil {
		fmt.Println("Erro ao enviar DM: ", err)
	}
}
