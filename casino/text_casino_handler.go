package casino

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var userStates = make(map[string]bool)

func formatBet(prefix string, ctx string) string {
	ctx = strings.TrimPrefix(ctx, prefix)
	ctx = strings.TrimSpace(ctx)
	return ctx
}

func ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch {
	case strings.HasPrefix(m.Content, "bj"):
		m.Content = formatBet("bj", m.Content)
		go startBlackJack(s, m, &userStates)
	// case strings.HasPrefix(m.Content, "!hilo"):
	// 	hiloGame(s, m)
	// case strings.HasPrefix(m.Content, "!dr"):
	// 	deathRoll(s, m)
	case strings.HasPrefix(m.Content, "daily"):
		player := m.Author.Username
		reply := dailyCoins(player)
		s.ChannelMessageSend(m.ChannelID, reply)
	default:
		return
	}
}

// func canPlay(authorID string, bet int) (bool, int) {
// 	isPlaying, ok := userStates[authorID]
// 	//check if user is playing if they are, my code is -1 to tell them they are already playing
// 	if ok && isPlaying {
// 		return false, -1
// 	}
// 	//if user is not in the states, add them
// 	if !ok {
// 		userStates[authorID] = false
// 	}
// 	bet = bet * 100
// 	//currently making the database file for balance check
// 	return true, 0 //0 needs to get changed to balance once logic is in place
// }
