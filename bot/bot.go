package bot

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var S *discordgo.Session
var Chans = make(map[string]chan *discordgo.InteractionCreate)

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set in the environment variables.")
	}

	S, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	S.Identify.Intents = discordgo.IntentsAll

	if S == nil {
		log.Fatal("Bot session is not initialized.")
	}

	err = S.Open()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is running!")

}
