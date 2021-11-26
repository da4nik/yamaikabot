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
	offset := 0
	for {
		resp, err := ta.tgClient.GetUpdates(telegram.GetUpdatesRequest{
			Timeout:        30,
			AllowedUpdates: []string{"message"},
			Offset:         offset,
		})
		if err != nil {
			ta.log.Errorf("Unable to get update: %s", err.Error())
			continue
		}

		if !resp.Ok {
			ta.log.Errorf("Request failed", resp)
			continue
		}

		ta.log.Debugf("Got update via LP: %+v", resp)
		for i := range resp.Updates {
			ta.handleUpdate(resp.Updates[i])
			offset = resp.Updates[i].UpdateId + 1
		}
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
		Ctx:        context.WithValue(context.Background(), "User", update.Message.From),
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
	user := msg.Ctx.Value("User").(telegram.User)

	if user.Id == 0 {
		ta.log.Errorf("Unable to send message, no user in context")
		return
	}

	ta.log.Debugf("Sending message to: %+v", user)
	_, err := ta.tgClient.SendMessage(telegram.SendMessageRequest{
		ChatId: user.Id,
		Text:   msg.Text,
	})
	if err != nil {
		ta.log.Errorf("Unable to send message to %s: %s", user.Username, err.Error())
	}
}
