package casino

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/casino/games"
)

func ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch {
	case strings.HasPrefix(m.Content, "!bj"):
		games.StartBlackJack(s, m)
	// case strings.HasPrefix(m.Content, "!hilo"):
	// 	hiloGame(s, m)
	// case strings.HasPrefix(m.Content, "!dr"):
	// 	deathRoll(s, m)
	default:
		// Optionally, handle unrecognized commands
		return
	}
}
