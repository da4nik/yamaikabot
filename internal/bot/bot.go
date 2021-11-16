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
	AnswerChan chan *Answer
	Ctx        context.Context
}

type Bot struct {
	In chan *Message

	log        logger.Logger
	processing bool
	done       chan bool
}

type botHandlerFunc func(string) (*Answer, error)

var commandRe = regexp.MustCompile(`^\/(\w+)\s*(.*)$`)

func New(queueSize int, log logger.Logger) *Bot {
	return &Bot{
		In:         make(chan *Message, queueSize),
		log:        log,
		processing: true,
		done:       make(chan bool),
	}
}

func (b *Bot) Start(ctx context.Context) {
	b.log.Infof("Starting ...")
	go func() {
		b.log.Infof("Started.")
		processing := true
		for processing {
			select {
			case <-ctx.Done():
				processing = false
			case msg := <-b.In:
				b.log.Debugf("Got message: %+v", msg)
				go b.processMessage(msg)
			}
		}
		b.log.Infof("Stopped.")
	}()
}

func (b *Bot) processMessage(msg *Message) {
	b.log.Debugf("Processing message %+v", msg)
	if msg.AnswerChan == nil {
		b.log.Errorf("Message won't be processed, answer chan is nil")
		return
	}

	cmd, rest, err := b.parseCommand(msg.Text)
	if err != nil {
		b.log.Errorf("Error parsing message: %s", err.Error())
		return
	}

	b.log.Infof("Command: %s, Rest: %s", cmd, rest)
	answer, err := b.handleCommand(cmd, rest)
	if err != nil {
		b.log.Errorf("Unable to process command `%s`: %s", cmd, err.Error())
		return
	}

	msg.AnswerChan <- answer
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

func (b *Bot) handleCommand(command, params string) (*Answer, error) {
	switch strings.ToLower(command) {
	case "start":
		return b.botNotImplemented(params)
	case "echo":
		return b.botEcho(params)
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}
