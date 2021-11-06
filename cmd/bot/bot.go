package main

import (
	"github.com/da4nik/yamaikabot/internal/bot"
	"github.com/da4nik/yamaikabot/internal/config"
	"github.com/da4nik/yamaikabot/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.ReadConfig()
	log := logger.NewLogger(conf.LogFile, conf.LogLevel)

	// Creating the bot
	theBot := bot.New(10, log.WithModule("bot"))

	// Run telegram bot routine
	// Run viber bot routine
	// Run bot itself
	theBot.Start()
	defer theBot.Stop()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}
