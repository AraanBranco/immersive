package roleplay

import (
	"bytes"
	"fmt"
)

func GetCities(ctx Context) {
	buffer := bytes.NewBufferString("Cidades disponíveis: \n")

	buffer.WriteString("```")
	for _, city := range ctx.Config.GetEnvConfStringSlice("cities") {
		msg := fmt.Sprintf("- %s \n", city)
		buffer.WriteString(msg)
	}
	buffer.WriteString("```")

	str := buffer.String()
	ctx.Reply(str)
}
