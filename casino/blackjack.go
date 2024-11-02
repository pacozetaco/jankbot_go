package casino

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"github.com/pacozetaco/jankbot_go/bot"
)

func startBlackJack(player string, mID string, bet int, bal int) {
	if bet%2 != 0 {
		_, err := bot.S.ChannelMessageSend(mID, "Bet must be even to play BlackJack.")
		if err != nil {
			log.Println(err)
		}
		userStates[player] = false
		return
	}

	game := &blackJackG{
		bG: bG{
			player:   player,
			bal:      bal,
			bet:      bet,
			mID:      mID,
			pAButton: pAButton,
		},
		hit:     hitButton,
		stand:   standButton,
		split:   splitButton,
		doubled: doubleDButton,
	}

	game.generateDeck(1)

	err := game.initializeBJ()
	if err != nil {
		log.Println(err)
		userStates[player] = false
		return
	}

	if game.choice != "blackjack" {
		game.playerBJTurn()
		game.jBBJTurn()
	}

	game.bJlogic()
	game.gameTransact()
	//if blackjack set the bet back to original value
	if game.choice == "blackjack" && game.result == "won" {
		game.bet = (game.bet / 3) * 2
	}

	game.logBJ()
	game.drawGame(false)

	if ch, ok := bot.Chans[game.board.ID]; ok {
		close(ch)
		delete(bot.Chans, game.board.ID)
		err = bot.S.ChannelMessageDelete(game.mID, game.board.ID)
		if err != nil {
			log.Println(err)
		}
	}

	game.sendComplex("", nil)

	os.Remove(game.pic)

	game.endGame(startBlackJack)
}

func (g *blackJackG) initializeBJ() error {
	for i := 0; i < 2; i++ {
		g.dealCard("player")
		g.dealCard("jb")
	}
	g.playerHandValue = bJHandValue(g.playerHand)
	g.jBHandValue = bJHandValue(g.jBHand)
	if g.playerHandValue == 21 || g.jBHandValue == 21 {
		g.choice = "blackjack"
	}
	content := "Shuffling Deck...."
	err := g.sendComplex(content, nil)
	if err != nil {
		return err
	}
	g.msg.Content = ""
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

func (g *blackJackG) bJlogic() {

	switch g.choice {
	case "bust":
		g.result = "lost"
		return
	case "blackjack":
		if g.playerHandValue == 21 {
			g.bet = g.bet + (g.bet / 2)
			g.result = "won"
		} else {
			g.result = "lost"
		}
		return
	case "timeout":
		g.result = "lost"
		return
	}
	switch {
	case g.jBHandValue > 21:
		g.result = "won"
		return
	case g.playerHandValue > g.jBHandValue:
		g.result = "won"
	case g.playerHandValue == g.jBHandValue:
		g.result = "pushed"
	default:
		g.result = "lost"
	}
}

func (g *blackJackG) jBBJTurn() {
	g.jBHandValue = bJHandValue(g.jBHand)
	if g.choice == "bust" || g.choice == "timeout" {
		return
	}
	for g.jBHandValue < 17 {
		g.dealCard("jb")
		g.jBHandValue = bJHandValue(g.jBHand)
	}
}

func (g *blackJackG) playerBJTurn() {
	g.jBHandValue = bJHandValue(g.jBHand[:1])
gameLoop:
	for {
		g.playerHandValue = bJHandValue(g.playerHand)
		close(bot.Chans[g.board.ID])
		delete(bot.Chans, g.board.ID)
		err := bot.S.ChannelMessageDelete(g.mID, g.board.ID)
		if err != nil {
			log.Println(err)
		}
		if g.playerHandValue > 21 {
			g.choice = "bust"
			break gameLoop
		}
		if g.playerHandValue == 21 {
			break gameLoop
		}

		g.drawGame(true)

		var buttons []discordgo.Button
		switch g.choice {
		case "":
			if g.bal >= g.bet*2 {
				buttons = []discordgo.Button{*g.hit, *g.stand, *g.doubled}
			} else {
				buttons = []discordgo.Button{*g.hit, *g.stand}
			}
		default:
			buttons = []discordgo.Button{*g.hit, *g.stand}
		}

		err = g.sendComplex("", buttons)
		if err != nil {
			log.Println(err)
		}

		err = g.handleButtonClick()
		if err != nil {
			log.Println(err)
			break gameLoop
		}

		switch g.choice {
		case "hit":
			g.dealCard("player")
		case "stand":
			break gameLoop
		case "double":
			g.bet = g.bet * 2
			g.dealCard("player")
			g.playerHandValue = bJHandValue(g.playerHand)
			if g.playerHandValue > 21 {
				g.choice = "bust"
			}
			break gameLoop
		case "timeout":
			break gameLoop
		}
	}
}

func (g *bG) pasteCards(dc *gg.Context, hand []string, y int) {
	cardWidth := 32
	totalWidth := len(hand) * cardWidth
	startX := (dc.Width() - totalWidth) / 2

	for i, card := range hand {
		cardPath := fmt.Sprintf("./assets/cards/%s.png", card)
		cardImage, _ := gg.LoadImage(cardPath)
		dc.DrawImage(cardImage, startX+i*cardWidth, y)
	}
}

func (g *blackJackG) drawGame(hide bool) error {
	const fontSize = 15
	fontPath := "./assets/font/pixel_font.ttf"
	tablePath := "./assets/tables/blackjack_table.png"
	gamePicPath := fmt.Sprintf("./temp/%s_game.png", g.player)

	// Determine dealer hand
	var dealerHand []string
	if hide {
		dealerHand = []string{"0_back", g.jBHand[0]}
	} else {
		dealerHand = g.jBHand
	}

	tableImage, err := gg.LoadImage(tablePath)
	if err != nil {
		log.Println(err)
		return err
	}

	dc := gg.NewContextForImage(tableImage)

	// Load font
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		log.Println(err)
		return err
	}
	// Draw cards
	g.pasteCards(dc, g.playerHand, 176)
	g.pasteCards(dc, dealerHand, 32)

	// Draw text
	dc.SetRGB(0, 0, 0)
	dc.DrawString(fmt.Sprintf("Jank: %d", g.jBHandValue), 23, 124)
	dc.DrawString(fmt.Sprintf("You: %d", g.playerHandValue), 38, 146)
	dc.DrawString("Bet", 176, 124)
	dc.DrawString(fmt.Sprintf("%d", g.bet), 176, 146)

	if g.result != "" {
		switch g.result {
		case "won":
			dc.DrawString("WINNER", 10, 250)
		case "lost":
			dc.DrawString("LOSER", 10, 250)
		case "pushed":
			dc.DrawString("PUSH", 10, 250)
		}
		dc.DrawString(fmt.Sprintf("Coins:%d", g.bal), 10, 20)
	}

	// Save image
	if err := dc.SavePNG(gamePicPath); err != nil {
		log.Println(err)
		return err
	}

	g.pic = gamePicPath
	return nil
}
