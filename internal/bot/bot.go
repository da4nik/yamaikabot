package bot

import (
	"context"
	"fmt"
	"github.com/da4nik/yamaikabot/internal/logger"
	"regexp"
	"strings"
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

	log        logger.Logger
	processing bool
	done       chan bool
}

var commandRe = regexp.MustCompile(`^\/(\w+)\s*(.*)$`)

func New(queueSize int, log logger.Logger) *Bot {
	return &Bot{
		In:         make(chan Message, queueSize),
		log:        log,
		processing: true,
		done:       make(chan bool),
	}
}

func (b *Bot) Start() {
	b.log.Infof("Starting ...")
	go func() {
		b.log.Infof("Started.")
		for b.processing {
			select {
			case <-b.done:
				b.processing = false
			case msg := <-b.In:
				b.log.Debugf("Got message: %+v", msg)
				go b.processMessage(msg)
			}
		}
		b.log.Infof("Stopping ...")
		b.done <- true
	}()
}

func (b *Bot) Stop() {
	b.done <- true
	<-b.done
	b.log.Infof("Stopped.")
}

func (b *Bot) processMessage(msg Message) {
	b.log.Debugf("Processing message %+v", msg)

	cmd, rest, err := b.parseCommand(msg.Text)
	if err != nil {
		b.log.Errorf("Error parsing message: %s", err.Error())
		return
	}

	b.log.Infof("Command: %s, Rest: %s", cmd, rest)
}

func (b *Bot) parseCommand(msg string) (string, string, error) {
	msg = strings.TrimSpace(msg)
	res := commandRe.FindAllStringSubmatch(msg, -1)

	b.log.Debugf("RES: %+v", res)

	if len(res) == 0 {
		return "", "", fmt.Errorf("command not found in string `%s`", msg)
	}

	matches := res[0]
	command := matches[1]
	rest := ""
	if len(matches) > 2 {
		rest = matches[2]
	}

	return command, rest, nil
}
