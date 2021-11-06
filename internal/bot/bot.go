package bot

import (
	"context"
	"github.com/da4nik/yamaikabot/internal/logger"
)

type Answer struct {
	Text string
	Ctx  context.Context
}

type Message struct {
	Text       string
	AnswerChan chan Answer
	Ctx        context.Context
}

type Bot struct {
	In chan Message

	log logger.Logger
}

func New(queueSize int, log logger.Logger) *Bot {
	return &Bot{
		In:  make(chan Message, queueSize),
		log: log,
	}
}

func (b *Bot) Start() {
	b.log.Infof("Starting ...")
}

func (b *Bot) Stop() {
	b.log.Infof("Stopping ...")
}
