package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitea.home/smythbrjr/jankbot_go/text_casino/casino"
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

	sess.AddHandler(onMessage)
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
		go casino.HandleMessage(s, m)
	case "ai-chat":
		println("we got an AI chat message")
		//route to AI chat module
	case "jukebox-spam":
		println("we got a jukebox spam message")
		//route to jukebox request
	case "ark-chat":
		println("we got an ARK chat message")
		//route to ARK module
	case "ark-config":
		println("we got an ARK config message")
		//route to arkconfig
	default:
		return
	}
}
