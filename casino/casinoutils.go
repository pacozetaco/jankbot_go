package casino

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

// play again button for all the games
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

// hilo struct
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
			clicker := i.Member.User.Username
			if clicker == g.player {
				err := bot.S.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredMessageUpdate,
				})
				if err != nil {
					fmt.Printf("Failed to respond to interaction: %v\n", err)
					return err
				}
				g.choice = i.MessageComponentData().CustomID
				return nil
			} else {
				err := bot.S.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredMessageUpdate,
				})
				if err != nil {
					fmt.Printf("Failed to respond to interaction: %v\n", err)
					return err
				}
			}
		case <-timer.C:
			fmt.Println("Timer expired, setting choice to timeout.")
			g.choice = "timeout"
			return nil
		}
	}
}

// updating gameboard func
func (g *bG) updateComplex() {
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

// add or subtract moneys from player depening on result
func (g *bG) gameTransact() {

	switch g.result {
	case "won":
		addBalance(g.player, g.bet)
		g.bal += g.bet
	case "lost":
		addBalance(g.player, -g.bet)
		g.bal -= g.bet
	}
}
