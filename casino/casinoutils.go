package casino

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

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

type diceG struct {
	roll int
}

type hiLoG struct {
	bG
	diceG
	hi *discordgo.Button
	lo *discordgo.Button
}

func (g *bG) handleButtonClick() {
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()

	for {
		select {
		case i := <-bot.Channels[g.board.ID]:
			clicker := i.Member.User.Username
			if clicker == g.player {
				fmt.Printf("User %s clicked button: %s\n", clicker, i.MessageComponentData().CustomID)
				err := bot.S.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredMessageUpdate,
				})
				if err != nil {
					fmt.Printf("Failed to respond to interaction: %v\n", err)
				}
				g.choice = i.MessageComponentData().CustomID
				fmt.Println(g.choice)
				return
			} else {
				err := bot.S.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredMessageUpdate,
				})
				if err != nil {
					fmt.Printf("Failed to respond to interaction: %v\n", err)
				}
			}
		case <-timer.C:
			fmt.Println("Timer expired, setting choice to timeout.")
			g.choice = "timeout"
			return
		}
	}
}

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
