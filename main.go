package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/aichat"
	"github.com/pacozetaco/jankbot_go/bot"
	"github.com/pacozetaco/jankbot_go/casino"
)

func main() {
	bot.StartBot()

	// Check if bot.S is initialized
	if bot.S == nil {
		log.Fatal("Bot session is not initialized.")
	}

	bot.S.AddHandler(onMessage)

	err := bot.S.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer bot.S.Close()

	log.Println("Bot is running!")
	casino.StartDb()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}

	switch channel.Name {
	case "casino":
		casino.ProcessCommand(m)
	case "ai-chat":
		go aichat.Chat(m)
	case "jukebox-spam":
		log.Println("Received a jukebox spam message")
	case "ark-chat":
		log.Println("Received an ARK chat message")
	case "ark-config":
		log.Println("Received an ARK config message")
	default:
		return
	}
}
