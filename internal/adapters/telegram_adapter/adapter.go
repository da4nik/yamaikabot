package telegram_adapter

import (
	"context"
	"fmt"

	"github.com/da4nik/yamaikabot/internal/bot"
	"github.com/da4nik/yamaikabot/internal/logger"
	"github.com/da4nik/yamaikabot/pkg/telegram"
)

type TelegramAdapter struct {
	log         logger.Logger
	longPooling bool
	webhookURL  string
	inputChan   chan *bot.Answer
	outputChan  chan *bot.Message
	tgClient    *telegram.Client
}

type TelegramAdapterOptions struct {
	Log         logger.Logger
	LongPooling bool
	WebhookURL  string
	BotToken    string
	OutputChan  chan *bot.Message
}

func New(opts TelegramAdapterOptions) (*TelegramAdapter, error) {
	if opts.BotToken == "" {
		return nil, fmt.Errorf("bot token is not provided")
	}

	if !opts.LongPooling && opts.WebhookURL == "" {
		return nil, fmt.Errorf("webhook url is not provided for webhook updates")
	}

	if opts.OutputChan == nil {
		return nil, fmt.Errorf("output chan is not provided")
	}

	return &TelegramAdapter{
		log:         opts.Log,
		longPooling: opts.LongPooling,
		webhookURL:  opts.WebhookURL,
		outputChan:  opts.OutputChan,
		inputChan:   make(chan *bot.Answer, 100),
		tgClient:    telegram.NewClient(opts.BotToken),
	}, nil
}

func (ta *TelegramAdapter) Start(ctx context.Context) {
	if ta.longPooling {
		go ta.startLongPooling(ctx)
	} else {
		go ta.startWebhooks(ctx)
	}

	go ta.botListener(ctx)

	go func(ctx context.Context) {
		<-ctx.Done()
	}(ctx)
}

func (ta *TelegramAdapter) startLongPooling(ctx context.Context) {
	ta.log.Infof("Starting long pooling telegram updates")
	for {
		update, err := ta.tgClient.GetUpdates(telegram.GetUpdatesRequest{
			Timeout:        30,
			AllowedUpdates: []string{"message"},
		})
		if err != nil {
			ta.log.Errorf("Unable to get update: %s", err.Error())
			continue
		}

		ta.log.Debugf("Got update via LP: %+v", update)
	}
}

func (ta *TelegramAdapter) startWebhooks(ctx context.Context) {
	err := ta.tgClient.SetWebhook(telegram.SetWebhookRequest{
		URL:                ta.webhookURL,
		DropPendingUpdates: true,
	})
	if err != nil {
		ta.log.Errorf("Unable to set webhook: %s", err.Error())
		return
	}

	defer func() {
		err := ta.tgClient.DeleteWebhook(false)
		if err != nil {
			ta.log.Errorf("Unable to remove webhook: %s", err.Error())
		}
	}()

	ta.log.Infof("Starting webhooks telegram updates")
}

func (ta *TelegramAdapter) handleUpdate(update telegram.Update) {
	ta.outputChan <- &bot.Message{
		Text:       update.Message.Text,
		Ctx:        context.Background(),
		AnswerChan: ta.inputChan,
	}
}

func (ta *TelegramAdapter) botListener(ctx context.Context) {
	for {
		select {
		case msg := <-ta.inputChan:
			ta.sendMsg(msg)
		case <-ctx.Done():
		}
	}
}

func (ta *TelegramAdapter) sendMsg(msg *bot.Answer) {

}
