package roleplay

import (
	"bytes"
	"fmt"
)

func HelpCommand(ctx Context) {
	cmds := ctx.CmdHandler.GetCmds()
	buffer := bytes.NewBufferString("Comandos: \n")
	for cmdName, cmdStruct := range cmds {
		if len(cmdName) == 1 {
			continue
		}

		msg := fmt.Sprintf("> %s%s - %s\n", ctx.Conf.Prefix, cmdName, cmdStruct.GetHelp())
		buffer.WriteString(msg)
	}
	str := buffer.String()
	ctx.Reply(str[:len(str)-2])
}