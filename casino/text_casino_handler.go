package casino

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

var (
	userStates = make(map[string]bool)
)

func ProcessCommand(m *discordgo.MessageCreate) {

	switch {
	case strings.HasPrefix(m.Content, "daily"):
		reply := dailyCoins(m.Author.Username)
		bot.S.ChannelMessageSend(m.ChannelID, reply)
		return
	}

	gameMap := map[string]func(player string, mID string, bet int, bal int){
		"hilo": startHiLo,
		"bj":   startBlackJack,
		"dr":   startDeathRoll,
	}

	for prefix, handler := range gameMap {
		if strings.HasPrefix(m.Content, prefix) {
			//trim up the prefix
			content := strings.TrimPrefix(m.Content, prefix)
			//clean up space
			content = strings.TrimSpace(content)
			//convert bet to int instead of string
			bet, err := strconv.Atoi(content)
			if err != nil {
				bot.S.ChannelMessageSend(m.ChannelID, "Invalid bet, check your syntax.")
				return
			}
			ok, bal, reply := canPlay(m.Author.Username, bet)
			{
				//check if user isnt in a game and has enough coins to play
				if !ok {
					bot.S.ChannelMessageSend(m.ChannelID, reply)
					//else start the game and add them to userStates
				} else {
					userStates[m.Author.Username] = true
					go handler(m.Author.Username, m.ChannelID, bet, bal)
				}
			}
			return
		}
	}
}

func canPlay(authorID string, bet int) (bool, int, string) {
	bal, err := getBalance(authorID)
	if err != nil {
		log.Println(err)
		bal = 0
	}
	isPlaying, ok := userStates[authorID]
	switch {
	case !ok:
		userStates[authorID] = false
	case isPlaying:
		return false, bal, "You're already playing. "
	}
	switch {
	case bal < bet:
		return false, bal, "You don't have enough coins. Balance: " + strconv.Itoa(bal)
	default:
		return true, bal, ""

	}
}
