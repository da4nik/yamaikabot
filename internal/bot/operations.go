package bot

import (
	"context"
	"fmt"
)

func (b *Bot) botNotImplemented(params string) (*Answer, error) {
	return nil, fmt.Errorf("not implemented")
}

func (b *Bot) botEcho(ctx context.Context, params string) (*Answer, error) {
	return &Answer{
		Text: fmt.Sprintf("ECHO: %s", params),
		Ctx:  ctx,
	}, nil
}
