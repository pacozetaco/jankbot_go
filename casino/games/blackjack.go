package games

import (
	"github.com/bwmarrin/discordgo"
)

func StartBlackJack(s *discordgo.Session, m *discordgo.MessageCreate) {

	println("Starting a new blackjack game!")

	_, err := s.ChannelMessageSend(m.ChannelID, "A new game of blackjack has started!")
	if err != nil {
		println(err)
	}
}
