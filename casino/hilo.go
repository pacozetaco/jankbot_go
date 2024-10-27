package casino

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

type HiLoGame struct {
	player string
	bal    int
	bet    int
	choice string
	roll   int
	result string
	board  *discordgo.Message
	mID    string
}

func startHiLo(player string, mID string, bet int, bal int) {

	game := HiLoGame{
		player: player,
		bal:    bal,
		bet:    bet,
		mID:    mID,
	}
	game.sendMessage()
}

func (h *HiLoGame) sendMessage() {
	var err error
	h.board, err = bot.S.ChannelMessageSend(h.mID, "Hi "+h.player+"! I am a game of Hi/Lo. You have "+strconv.Itoa(h.bal)+" coins. What is your bet?")
	if err != nil {
		log.Println(err)
		return
	}
}
