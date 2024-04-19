package tg

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"log/slog"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/handler/app"
)

type Handler struct {
	log *slog.Logger
	bot bot.Bot
}

func New(log *slog.Logger, bot bot.Bot) *Handler {
	return &Handler{
		log: log,
		bot: bot,
	}
}

func (h *Handler) Handle() {
	h.bot.Bh.Handle(func(bot *telego.Bot, update telego.Update) {
		_, _ = bot.SendMessage(
			tu.Messagef(
				tu.ID(update.Message.Chat.ID),
				app.MsgStart, update.Message.From.Username,
			).WithParseMode(telego.ModeMarkdownV2),
		)
	}, th.CommandEqual(app.StartCmd.Command))
}

func (h *Handler) HandleEnd() {
	h.bot.Bh.Handle(func(bot *telego.Bot, update telego.Update) {
		_, _ = bot.SendMessage(
			tu.Messagef(
				tu.ID(update.Message.Chat.ID),
				app.MsgStart, update.Message.From.Username,
			).WithParseMode(telego.ModeMarkdownV2),
		)
	}, th.Not(th.AnyCommand()))
}
