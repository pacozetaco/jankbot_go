package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}
	})
	sess.Identify.Intents = discordgo.IntentsAll

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	println("Bot is running!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
