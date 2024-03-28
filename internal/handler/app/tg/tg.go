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
	h.bot.Bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		h.log.Info("got new command", message.From.Username, message.Text)

		_, _ = bot.SendMessage(
			tu.Messagef(
				tu.ID(message.Chat.ID),
				app.MsgStart, message.From.Username,
			).WithParseMode(telego.ModeMarkdownV2),
		)
	}, th.CommandEqual(app.StartCmd.Command))
}
