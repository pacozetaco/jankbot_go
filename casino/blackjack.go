package casino

import (
	"log"
	"strings"
)

func startBlackJack(player string, mID string, bet int, bal int) {

	game := &blackJackG{
		bG: bG{
			player:   player,
			bal:      bal,
			bet:      bet,
			mID:      mID,
			pAButton: pAButton,
		},
		whosturn: "",
		hit:      hitButton,
		stand:    standButton,
		split:    splitButton,
		doubled:  doubleDButton,
	}

	game.generateDeck(1)

	err := game.initializeBJ()
	if err != nil {
		log.Println(err)
		userStates[player] = false
		return
	}
}

func (g *blackJackG) initializeBJ() error {
	for i := 0; i < 2; i++ {
		g.dealCard("player")
		g.dealCard("jb")
	}
	g.playerHandValue = bJHandValue(g.playerHand)
	g.jBHandValue = bJHandValue(g.jBHand)
	if g.playerHandValue == 21 || g.jBHandValue == 21 {
		g.whosturn = "blackjack"
	}
	content := "Shuffling Deck...."
	err := g.sendComplex(content, nil)
	if err != nil {
		return err
	}
	return nil
}

func bJHandValue(hand []string) int {
	cardValues := map[string]int{
		"0":  0,
		"2":  2,
		"3":  3,
		"4":  4,
		"5":  5,
		"6":  6,
		"7":  7,
		"8":  8,
		"9":  9,
		"10": 10,
		"j":  10,
		"q":  10,
		"k":  10,
		"a":  11,
	}
	var aces int
	var handvalue int
	for _, card := range hand {
		cardParts := strings.Split(card, "_")
		cardValue := cardValues[string(cardParts[0])]
		handvalue += cardValue
		if string(cardParts[0]) == "a" {
			aces++
		}
	}
	for handvalue > 21 && aces > 0 {
		handvalue -= 10
		aces--
	}
	return handvalue

}
