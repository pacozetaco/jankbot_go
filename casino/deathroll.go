package casino

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

func startDeathRoll(m *discordgo.MessageCreate) {
	bot.S.ChannelMessageSend(m.ChannelID, "Deathroll")
}
