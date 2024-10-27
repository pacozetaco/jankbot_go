package casino

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

func startHiLo(player string, chanID string, bet int, bal int) {

	hiButton := &discordgo.Button{
		Label:    "Hi",
		Style:    3,
		Disabled: false,
		CustomID: "high",
	}

	loButton := &discordgo.Button{
		Label:    "Lo",
		Style:    4,
		Disabled: false,
		CustomID: "low",
	}

	game := hiLoG{
		bG: bG{
			player:   player,
			bal:      bal,
			bet:      bet,
			mID:      chanID,
			pAButton: pAButton,
		},
		diceG: diceG{
			roll: rand.Intn(100) + 1,
		},
		hi: hiButton,
		lo: loButton,
	}

	err := game.initializeHiLo()
	if err != nil {
		log.Println(err)
		userStates[player] = false
		return
	}
	err = game.handleButtonClick()
	if err != nil {
		log.Println(err)
		userStates[player] = false
		return
	}
	game.gameLogic()
	game.endHiLo()

}

func (h *hiLoG) initializeHiLo() error {
	var err error
	content := fmt.Sprintf("```HiLo! Bet: %d\nIs your roll higher or lower than 50? (/roll 1-100)```", h.bet)
	h.msg = &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{h.hi, h.lo},
			},
		},
	}
	h.board, err = bot.S.ChannelMessageSendComplex(h.mID, h.msg)
	if err != nil {
		log.Println(err)
		userStates[h.player] = false
		return err
	}
	bot.Chans[h.board.ID] = make(chan *discordgo.InteractionCreate)
	return nil
}

func (h *hiLoG) gameLogic() {
	switch h.choice {
	case "high":
		if h.roll > 50 {
			h.result = "won"
		} else if h.roll < 50 {
			h.result = "lost"
		} else {
			h.result = "tie"
		}
	case "low":
		if h.roll > 50 {
			h.result = "lost"
		} else if h.roll < 50 {
			h.result = "won"
		} else {
			h.result = "tie"
		}
	case "timeout":
		h.result = "lost"
	}
	h.gameTransact()
}

func (h *hiLoG) endHiLo() {
	content := fmt.Sprintf("\n```You chose %s and rolled a %d\nResult: %s  Balance: %d```", h.choice, h.roll, h.result, h.bal)
	h.msg.Content += content

	if h.bal > h.bet {
		h.msg.Components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{h.pAButton},
			},
		}
	} else {
		h.msg.Components = []discordgo.MessageComponent{}
	}
	h.updateComplex()

	userStates[h.player] = false

	h.logHiLo()

	if h.msg.Components != nil {
		h.handleButtonClick()
		h.msg.Components = []discordgo.MessageComponent{}
		h.updateComplex()
	}
	close(bot.Chans[h.board.ID])
	delete(bot.Chans, h.board.ID)

	if h.choice == "play" && h.bal > h.bet && !userStates[h.player] {
		startHiLo(h.player, h.mID, h.bet, h.bal)
	}
}
