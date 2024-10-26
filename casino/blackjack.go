package casino

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func startBlackJack(s *discordgo.Session, m *discordgo.MessageCreate, userStates *map[string]bool) {
	bet, err := strconv.Atoi(m.Content)
	if err != nil {
		log.Println(err)
		return
	}
	balance, err := getBalance(m.Author.Username)
	if err != nil {
		log.Println(err)
		return
	}

	if balance < bet {
		s.ChannelMessageSend(m.ChannelID, "You don't have enough coins")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Blackjack")
}
