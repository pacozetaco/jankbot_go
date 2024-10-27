package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pacozetaco/jankbot_go/bot"
	"github.com/pacozetaco/jankbot_go/handlers"
)

func main() {

	bot.S.AddHandler(handlers.OnMessage)
	bot.S.AddHandler(handlers.ButtonHandler)
	defer bot.S.Close()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
