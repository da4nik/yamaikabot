package listener

import (
	"context"
	"fmt"
	"github.com/da4nik/yamaikabot/internal/logger"
)

type Listener struct {
	messengers []ChatAdapter
	in         chan string
	log        logger.Logger
}

type ChatAdapter interface {
	Start(context.Context)
}

func New(messengers []ChatAdapter, log logger.Logger) (*Listener, error) {
	if len(messengers) == 0 {
		return nil, fmt.Errorf("no messengers provided")
	}

	return &Listener{
		log:        log,
		in:         make(chan string, 10),
		messengers: messengers,
	}, nil
}

func (l *Listener) Start(ctx context.Context) {
	for _, messenger := range l.messengers {
		messenger.Start(ctx)
	}

	go func() {
		<-ctx.Done()
	}()
}
