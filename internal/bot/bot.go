package bot

import (
	"github.com/mymmrac/telego"
)

type Bot struct {
	Bot     *telego.Bot
	Channel <-chan telego.Update
}

func New(token string) Bot {
	bot, err := telego.NewBot(token)
	if err != nil {
		panic("can't create bot: " + err.Error())
	}

	channel, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		panic("can't create bot's updates channel: " + err.Error())
	}

	return Bot{
		Bot:     bot,
		Channel: channel,
	}
}
