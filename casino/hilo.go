package casino

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func startHiLo(player string, chanID string, bet int, bal int) {

	game := &hiLoG{
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
	game.hilogameLogic()
	game.gameTransact()
	game.logHiLo()
	game.msg.Content += fmt.Sprintf("\nYou chose %s and rolled a %d. You %s %d coins.\nBalance: %d", game.choice, game.roll, game.result, game.bet, game.bal)
	game.endGame(startHiLo)

}

func (h *hiLoG) initializeHiLo() error {
	content := fmt.Sprintf("HiLo! Bet: %d\nIs your roll higher or lower than 50? (/roll 1-100)", h.bet)
	err := h.sendComplex(content, []discordgo.Button{*h.hi, *h.lo})
	if err != nil {
		return err
	}
	return nil
}

func (h *hiLoG) hilogameLogic() {
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
}
