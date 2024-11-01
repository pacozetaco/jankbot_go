package aichat

import (
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/JexSrs/go-ollama"
	"github.com/bwmarrin/discordgo"
	"github.com/pacozetaco/jankbot_go/bot"
)

func Chat(m *discordgo.MessageCreate) {
	preface := "You are a Discord bot. Your name is JankBot. " + m.Author.Username + " is the person you're talking to. You will do you best to answer any questions they may have. They say: "
	u, err := url.Parse("http://192.168.1.2:11434")
	if err != nil {
		log.Println(err)
		return
	}
	seed, err := strconv.Atoi(m.Message.ID)
	if err != nil {
		log.Println(err)
		return
	}
	llm := ollama.New(*u)
	res, err := llm.Generate(
		llm.Generate.WithModel("llama3.2:3b"),
		llm.Generate.WithPrompt(preface+m.Content),
		llm.Generate.WithSeed(seed),
	)

	if err != nil {
		log.Println(err)
		bot.S.ChannelMessageSend(m.ChannelID, "Sum-Ting-Wong")
		return
	}
	text := res.Response

	if len(text) > 1980 {
		words := strings.Split(text, " ")
		var chunks []string
		var currentChunk string

		for _, word := range words {
			if len(currentChunk)+len(word)+1 > 1980 {
				chunks = append(chunks, currentChunk+"...")
				currentChunk = word
			} else {
				if currentChunk != "" {
					currentChunk += " "
				}
				currentChunk += word
			}
		}

		if currentChunk != "" {
			chunks = append(chunks, currentChunk)
		}

		for _, chunk := range chunks {
			bot.S.ChannelMessageSend(m.ChannelID, "```"+chunk+"```")
		}
	} else {
		bot.S.ChannelMessageSend(m.ChannelID, "```"+text+"```")
	}
}
