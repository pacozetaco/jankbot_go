package casino

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

func startBlackJack(m *discordgo.MessageCreate) {
	bet, err := strconv.Atoi(m.Content)
	if err != nil {
		log.Println(err)
		bot.S.ChannelMessageSend(m.ChannelID, "Invalid bet")
		return
	}
	ok, bal, rply := canPlay(m.Author.Username, bet)

	if !ok {
		bot.S.ChannelMessageSend(m.ChannelID, rply+"Balance: "+strconv.Itoa(bal))
		return
	}

	(userStates)[m.Author.Username] = true
	bot.S.ChannelMessageSend(m.ChannelID, "Blackjack")
}
