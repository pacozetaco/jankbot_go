package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/aichat"
	"github.com/pacozetaco/jankbot_go/bot"
	"github.com/pacozetaco/jankbot_go/casino"
)

func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	println("we have a message")
	switch channel.Name {
	case "casino":
		go casino.ProcessCommand(m)
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

func ButtonHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if channel, ok := bot.Channels[i.Message.ID]; ok {
		channel <- i
	}
}
