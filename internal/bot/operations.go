package bot

import "fmt"

func (b *Bot) botNotImplemented(params string) (*Answer, error) {
	return nil, fmt.Errorf("not implemented")
}

func (b *Bot) botEcho(params string) (*Answer, error) {
	return &Answer{
		Text: fmt.Sprintf("ECHO: %s", params),
	}, nil
}
