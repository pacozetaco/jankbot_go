package casino

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func processCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch {
	case strings.HasPrefix(m.Content, "!bj"):
		println("starting a bj game")
		bjGame(s, m)
	case strings.HasPrefix(m.Content, "!hilo"):
		hiloGame(s, m)
	case strings.HasPrefix(m.Content, "!dr"):
		deathRoll(s, m)
	default:
		// Optionally, handle unrecognized commands
		return
	}
}
