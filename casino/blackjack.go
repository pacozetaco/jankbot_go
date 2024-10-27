package casino

import (
	"github.com/pacozetaco/jankbot_go/bot"
)

func startBlackJack(layer string, mID string, bet int, bal int) {

	bot.S.ChannelMessageSend(mID, "Blackjack")
}
