package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pacozetaco/jankbot_go/bot"
	"github.com/pacozetaco/jankbot_go/casino"
	"github.com/pacozetaco/jankbot_go/handlers"
)

func main() {

	bot.StartBot()

	casino.StartDb()

	bot.S.AddHandler(handlers.OnMessage)
	bot.S.AddHandler(handlers.ButtonHandler)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
