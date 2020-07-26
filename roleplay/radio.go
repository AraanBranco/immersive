package roleplay

import (
	"fmt"
	"math/rand"
	"time"
)

func RadioCommand(ctx Context) {
	rand.Seed(time.Now().UnixNano())
	min := 100
	max := 999

	rasult := rand.Intn(max-min) + min
	ctx.Reply(fmt.Sprintf("Frequência do rádio: **%v**", rasult))
}
