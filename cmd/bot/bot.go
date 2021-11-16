package main

import (
	"context"
	"fmt"
	tga "github.com/da4nik/yamaikabot/internal/adapters/telegram_adapter"
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
	ctx, cancel := context.WithCancel(context.Background())

	// Creating the bot
	theBot := bot.New(10, log.WithModule("bot"))

	tg, err := tga.New(tga.TelegramAdapterOptions{
		Log:         log.WithModule("telegram_adapter"),
		LongPooling: true,
		BotToken:    conf.TelegramBotToken,
		WebhookURL:  "https://makstep.ru/",
		OutputChan:  theBot.In,
	})
	if err != nil {
		log.Errorf("Unable to start telegram adapter: %s", err.Error())
		os.Exit(1)
	}

	// Run telegram bot routine
	tg.Start(ctx)

	// Run bot itself
	theBot.Start(ctx)

	answerChan := make(chan *bot.Answer)
	// theBot.In <- bot.Message{
	// 	Text:       "/start some name",
	// 	AnswerChan: answerChan,
	// }
	// fmt.Printf("start %+v\n", <-answerChan)

	theBot.In <- &bot.Message{
		Text:       "/echo hello there !!!",
		AnswerChan: answerChan,
	}
	fmt.Printf("echo %+v\n", <-answerChan)

	// theBot.In <- bot.Message{
	// 	Text:       "/unknown_command",
	// 	AnswerChan: answerChan,
	// }

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Infof("interrupt signal")
	cancel()
}
