package bot

import (
	"github.com/mymmrac/telego"
	"work-routine-bot/internal/processor/working-hours"
)

type Bot struct {
	Bot     *telego.Bot
	Channel <-chan telego.Update
}

var MenuCmds = []telego.BotCommand{
	working_hours.StartCmd,
	working_hours.ListWorkingHours,
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

	err = bot.SetMyCommands(&telego.SetMyCommandsParams{Commands: MenuCmds})
	if err != nil {
		panic("can't set bot's commands: " + err.Error())
	}

	return Bot{
		Bot:     bot,
		Channel: channel,
	}
}
