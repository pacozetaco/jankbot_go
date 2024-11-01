package casino

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

// BUTTONS
var pAButton = &discordgo.Button{
	Label:    "Play Again?",
	Style:    3,
	Disabled: false,
	CustomID: "play",
}

var rollbtn = &discordgo.Button{
	Label:    "Roll",
	Style:    1,
	Disabled: false,
	CustomID: "roll",
}

var autobtn = &discordgo.Button{
	Label:    "Auto Roll",
	Style:    3,
	Disabled: false,
	CustomID: "auto",
}

var youbutn = &discordgo.Button{
	Label:    "You",
	Style:    3,
	Disabled: false,
	CustomID: "you",
}

var jbbutn = &discordgo.Button{
	Label:    "JB",
	Style:    4,
	Disabled: false,
	CustomID: "jb",
}

var hiButton = &discordgo.Button{
	Label:    "Hi",
	Style:    3,
	Disabled: false,
	CustomID: "high",
}

var loButton = &discordgo.Button{
	Label:    "Lo",
	Style:    4,
	Disabled: false,
	CustomID: "low",
}

var hitButton = &discordgo.Button{
	Label:    "Hit",
	Style:    1,
	Disabled: false,
	CustomID: "hit",
}

var standButton = &discordgo.Button{
	Label:    "Stand",
	Style:    4,
	Disabled: false,
	CustomID: "stand",
}

var doubleDButton = &discordgo.Button{
	Label:    "Double Down",
	Style:    3,
	Disabled: false,
	CustomID: "double",
}

var splitButton = &discordgo.Button{
	Label:    "Split",
	Style:    4,
	Disabled: false,
	CustomID: "split",
}

// STRUCTS
// basic game stuct all games will use
type bG struct {
	player   string
	bal      int
	bet      int
	result   string
	board    *discordgo.Message
	msg      *discordgo.MessageSend
	mID      string
	choice   string
	pAButton *discordgo.Button
	pic      string
}

// card game struct
type cardG struct {
	deck       []string
	playerHand []string
	jBHand     []string
}

// dice game struc
type diceG struct {
	roll int
}

type deathRollG struct {
	bG
	diceG
	you       *discordgo.Button
	jb        *discordgo.Button
	rollbuttn *discordgo.Button
	autoroll  *discordgo.Button
	turn      string
	first     string
}

type blackJackG struct {
	bG
	cardG
	stand           *discordgo.Button
	hit             *discordgo.Button
	doubled         *discordgo.Button
	split           *discordgo.Button
	playerHandValue int
	jBHandValue     int
}

type hiLoG struct {
	bG
	diceG
	hi *discordgo.Button
	lo *discordgo.Button
}

// handles input from buttons, if its not the user it will drop the chan msg and look for the player.
// after 30 secs it will timeout and the player will lose their bet
func (g *bG) handleButtonClick() error {
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()

	for {
		select {
		case i := <-bot.Chans[g.board.ID]:
			err := bot.S.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})
			if err != nil {
				log.Printf("Failed to respond to interaction: %v\n", err)
				return err
			}
			clicker := i.Member.User.Username
			if clicker == g.player {
				g.choice = i.MessageComponentData().CustomID
				return nil
			}
		case <-timer.C:
			log.Println("Timer expired, setting choice to timeout.")
			g.choice = "timeout"
			return nil
		}
	}
}

// updating gameboard func
func (g *bG) updateComplex(buttons []discordgo.Button) {
	var components []discordgo.MessageComponent
	for _, button := range buttons {
		components = append(components, button) // Append each button as a MessageComponent
	}
	g.msg.Components = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: components,
		},
	}
	if buttons == nil {
		g.msg.Components = []discordgo.MessageComponent{}
	}
	updatedMsg := &discordgo.MessageEdit{
		Channel:    g.board.ChannelID,
		ID:         g.board.ID,
		Content:    &g.msg.Content,
		Components: &g.msg.Components,
	}

	_, err := bot.S.ChannelMessageEditComplex(updatedMsg)
	if err != nil {
		log.Println(err)
	}
}

func (b *bG) sendComplex(content string, buttons []discordgo.Button) error {
	// Create a slice of MessageComponent to hold the buttons
	var components []discordgo.MessageComponent
	var discordFile *discordgo.File

	file, err := os.Open(b.pic)
	if err == nil {
		discordFile = &discordgo.File{
			Name:   "game_image.png",
			Reader: file,
		}
	}
	defer file.Close()
	for _, button := range buttons {
		components = append(components, button) // Append each button as a MessageComponent
	}
	b.msg = &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: components,
			},
		},
		File: discordFile,
	}

	if buttons == nil {
		b.msg.Components = []discordgo.MessageComponent{}
	}

	b.board, err = bot.S.ChannelMessageSendComplex(b.mID, b.msg)
	if err != nil {
		log.Println(err)
		return err
	}

	bot.Chans[b.board.ID] = make(chan *discordgo.InteractionCreate, 10)
	return nil
}

type startFuncType func(playerID string, mID string, bet int, balance int)

func (g *bG) endGame(startFunc startFuncType) {
	if g.bal >= g.bet {
		g.updateComplex([]discordgo.Button{*pAButton})
	} else {
		g.updateComplex(nil)
	}

	userStates[g.player] = false

	if g.msg.Components != nil {
		g.handleButtonClick()
		g.updateComplex(nil)
	}
	close(bot.Chans[g.board.ID])
	delete(bot.Chans, g.board.ID)

	if g.choice == "play" && g.bal > g.bet && !userStates[g.player] {
		startFunc(g.player, g.mID, g.bet, g.bal)
	}
}

// add or subtract moneys from player depening on result
func (g *bG) gameTransact() {
	switch g.result {
	case "won":
		go addBalance(g.player, g.bet)
		g.bal += g.bet
	case "lost":
		go addBalance(g.player, -g.bet)
		g.bal -= g.bet
	}
}

func (c *cardG) generateDeck(numofDecks int) {
	cards := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "j", "q", "k", "a"}
	suits := []string{"spade", "heart", "diamond", "club"}
	for j := 0; j < numofDecks; j++ {
		for _, suit := range suits {
			for _, card := range cards {
				c.deck = append(c.deck, card+"_"+suit)
			}
		}
	}
	for i := 0; i < 10; i++ {
		rand.Shuffle(len(c.deck), func(i, j int) { c.deck[i], c.deck[j] = c.deck[j], c.deck[i] })
	}
}

func (c *cardG) dealCard(hand string) {
	if hand == "player" {
		c.playerHand = append(c.playerHand, c.deck[0])
		c.deck = c.deck[1:]
	} else if hand == "jb" {
		c.jBHand = append(c.jBHand, c.deck[0])
		c.deck = c.deck[1:]
	}
}
