package tg

import (
	"context"
	"github.com/mnogokotin/golang-packages/utils/e"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log/slog"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/storage/pages"
)

type Processor struct {
	log     *slog.Logger
	bot     bot.Bot
	storage pages.Storage
}

func New(log *slog.Logger, bot bot.Bot, storage pages.Storage) *Processor {
	return &Processor{
		log:     log,
		bot:     bot,
		storage: storage,
	}
}

func (p *Processor) GetChan() <-chan telego.Update {
	return p.bot.Channel
}

func (p *Processor) Process(ctx context.Context, u telego.Update) error {
	if err := p.executeCmd(ctx, u); err != nil {
		return e.Wrap("can't process update", err)
	}

	return nil
}

func (p *Processor) SendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := p.bot.Bot.SendMessage(
		tu.Message(
			tu.ID(chatID),
			text,
		),
	)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}
