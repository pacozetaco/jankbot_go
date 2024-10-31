package casino

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func startDeathRoll(player string, mID string, bet int, bal int) {

	game := &deathRollG{
		bG: bG{
			player:   player,
			bal:      bal,
			bet:      bet,
			mID:      mID,
			pAButton: pAButton,
		},
		diceG: diceG{
			roll: 100,
		},
		rollbuttn: rollbtn,
		you:       youbutn,
		jb:        jbbutn,
		autoroll:  autobtn,
	}

	err := game.initializeDR()
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
	err = game.drLoop()
	if err != nil {
		log.Println(err)
		userStates[player] = false
		return
	}
	game.drLogic()
	game.endGame(startDeathRoll)

}

func (d *deathRollG) initializeDR() error {
	content := fmt.Sprintf("Deathroll! Bet: %d\nWho goes first? (/roll 1-100)", d.bet)
	err := d.sendComplex(content, []discordgo.Button{*d.you, *d.jb})
	if err != nil {
		return err
	}
	return nil
}

func (d *deathRollG) drLoop() error {

	if d.choice == "timeout" {
		d.first = "Timeout"
		return nil
	}
	if d.turn == "" {
		switch d.choice {
		case "you":
			d.turn = "you"
			d.first = "Player"
			d.msg.Content += "\nYou go first"
		case "jb":
			d.turn = "jb"
			d.first = "Jankbot"
			d.msg.Content += "\nJankbot goes first"
		}
	}
	for {
		switch {
		case d.turn == "jb":
			d.roll = rand.Intn(d.roll) + 1
			d.msg.Content += fmt.Sprintf("\nJB rolled a %d", d.roll)
		case d.turn == "you" && d.choice != "auto":
			d.updateComplex([]discordgo.Button{*d.rollbuttn, *d.autoroll})
			err := d.handleButtonClick()
			if err != nil {
				log.Println(err)
				return err
			}
			switch d.choice {
			case "timeout":
				return nil
			case "roll":
				d.roll = rand.Intn(d.roll) + 1
				d.msg.Content += fmt.Sprintf("\nYou rolled a %d", d.roll)
			}
		}
		if d.turn == "you" && d.choice == "auto" {
			d.roll = rand.Intn(d.roll) + 1
			d.msg.Content += fmt.Sprintf("\nYou rolled a %d", d.roll)
		}
		if d.roll == 1 {
			return nil
		}
		switch d.turn {
		case "jb":
			d.turn = "you"
		case "you":
			d.turn = "jb"
		}
	}
}

func (d *deathRollG) drLogic() {
	switch d.turn {
	case "jb":
		d.result = "won"
	case "you":
		d.result = "lost"
	}
	d.gameTransact()
	d.logDeathRoll()
}
