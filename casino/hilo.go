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
	msg    *discordgo.MessageCreate
}

func startHiLo(m *discordgo.MessageCreate) {
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
	// Create a new instance of HiLoGame
	game := HiLoGame{
		player: m.Author.Username,
		bal:    bal,
		bet:    bet,
		msg:    m,
	}
	game.sendMessage()
}

func (h *HiLoGame) sendMessage() {
	var err error
	h.board, err = bot.S.ChannelMessageSend(h.msg.ChannelID, "Hi "+h.player+"! I am a game of Hi/Lo. You have "+strconv.Itoa(h.bal)+" coins. What is your bet?")
	if err != nil {
		log.Println(err)
		return
	}
}
