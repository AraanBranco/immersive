package radio

import (
	"math/rand"
	"time"
	"fmt"

	"github.com/bwmarrin/discordgo"

)

func Generate(s *discordgo.Session, m *discordgo.MessageCreate) {
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999

	rasult := rand.Intn(max-min)+min
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```CSS Rádio: %v ```", rasult))
}