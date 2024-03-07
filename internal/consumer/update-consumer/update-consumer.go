package update_consumer

import (
	"context"
	"github.com/mnogokotin/golang-packages/utils/e"
	"github.com/mymmrac/telego"
	"log/slog"
	"work-routine-bot/internal/processor"
)

type Consumer struct {
	log       *slog.Logger
	fetcher   processor.Fetcher
	processor processor.Processor
}

func New(log *slog.Logger, fetcher processor.Fetcher, processor processor.Processor) Consumer {
	return Consumer{
		log:       log,
		fetcher:   fetcher,
		processor: processor,
	}
}

func (c *Consumer) Start() error {
	updatesChan := c.fetcher.GetChan()

	for update := range updatesChan {
		if err := c.handleUpdate(context.Background(), update); err != nil {
			c.log.Error("can't handle updates", "", err.Error())

			continue
		}
	}

	return nil
}

func (c *Consumer) handleUpdate(ctx context.Context, update telego.Update) error {
	if err := c.processor.Process(ctx, update); err != nil {
		return e.Wrap("can't handle update", err)
	}

	return nil
}
