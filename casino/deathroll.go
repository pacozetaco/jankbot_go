package casino

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

func startDeathRoll(player string, mID string, bet int, bal int) {

	game := deathRollG{
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
	game.endDR()

}

func (d *deathRollG) initializeDR() error {
	content := fmt.Sprintf("Deathroll! Bet: %d\nWho goes first? (/roll 1-100)", d.bet)
	d.msg = &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{d.you, d.jb},
			},
		},
	}
	var err error
	d.board, err = bot.S.ChannelMessageSendComplex(d.mID, d.msg)
	if err != nil {
		log.Println(err)
		return err
	}
	bot.Chans[d.board.ID] = make(chan *discordgo.InteractionCreate)
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
			d.msg.Components = []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{d.rollbuttn, d.autoroll},
				},
			}
			d.updateComplex()
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
}

func (d *deathRollG) endDR() {
	if d.bal > d.bet {
		d.msg.Components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{d.pAButton},
			},
		}
	} else {
		d.msg.Components = []discordgo.MessageComponent{}
	}
	d.msg.Content += fmt.Sprintf("\nYou %s %d coins.\nBalance: %d", d.result, d.bet, d.bal)
	userStates[d.player] = false
	d.updateComplex()
	d.logDeathRoll()
	if d.msg.Components != nil {
		d.handleButtonClick()
		d.msg.Components = []discordgo.MessageComponent{}
		d.updateComplex()
	}
	close(bot.Chans[d.board.ID])
	delete(bot.Chans, d.board.ID)

	if d.choice == "play" && d.bal > d.bet && !userStates[d.player] {
		startDeathRoll(d.player, d.mID, d.bet, d.bal)
	}
}
