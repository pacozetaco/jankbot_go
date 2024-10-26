package casino

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

var (
	userStates = make(map[string]bool)
)

func formatBet(prefix string, ctx string) string {
	ctx = strings.TrimPrefix(ctx, prefix)
	ctx = strings.TrimSpace(ctx)
	return ctx
}

func ProcessCommand(m *discordgo.MessageCreate) {
	switch {
	case strings.HasPrefix(m.Content, "bj"):
		m.Content = formatBet("bj", m.Content)
		go startBlackJack(m)
	case strings.HasPrefix(m.Content, "hilo"):
		m.Content = formatBet("hilo", m.Content)
		go startHiLo(m)
	case strings.HasPrefix(m.Content, "dr"):
		m.Content = formatBet("dr", m.Content)
		go startDeathRoll(m)
	case strings.HasPrefix(m.Content, "daily"):
		player := m.Author.Username
		reply := dailyCoins(player)
		bot.S.ChannelMessageSend(m.ChannelID, reply)
	default:
		return
	}
}

func canPlay(authorID string, bet int) (bool, int, string) {
	bal, err := getBalance(authorID)
	if err != nil {
		log.Println(err)
		bal = 0
	}
	isPlaying, ok := userStates[authorID]
	switch {
	case !ok:
		userStates[authorID] = false
	case isPlaying:
		return false, bal, "You're already playing. "
	}
	switch {
	case bal < bet:
		return false, bal, "You don't have enough coins"
	default:
		return true, bal, ""

	}
}
