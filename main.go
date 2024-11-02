package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pacozetaco/jankbot_go/bot"
	"github.com/pacozetaco/jankbot_go/handlers"
)

func main() {
	//casino db self inits
	//bot self inits
	//start up the on msg and button handler
	bot.S.AddHandler(handlers.OnMessage)
	bot.S.AddHandler(handlers.ButtonHandler)
	// jankservers.StartServerMonitor()
	//close the bot connection on exit / failure
	defer bot.S.Close()
	//wait for system interupt to stop the bot
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
