package bot

import (
	"github.com/looplab/fsm"
	t "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"work-routine-bot/internal/handler/app"
	"work-routine-bot/internal/handler/task"
)

type Fsm struct {
	Fsm  *fsm.FSM
	Data interface{}
}

type Bot struct {
	Bot *t.Bot
	Bh  *th.BotHandler
	Fsm *Fsm
}

var MenuCmds = []t.BotCommand{
	app.StartCmd,
	task.MyTasks,
}

func New(token string, env string) Bot {
	bot, err := t.NewBot(token)
	if err != nil {
		panic("can't create bot: " + err.Error())
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		panic("can't create bot's updates channel: " + err.Error())
	}

	err = bot.SetMyCommands(&t.SetMyCommandsParams{Commands: MenuCmds})
	if err != nil {
		panic("can't set bot's commands: " + err.Error())
	}

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		panic("can't create bot's handler: " + err.Error())
	}

	fsm_ := &Fsm{Fsm: nil, Data: nil}

	return Bot{
		Bot: bot,
		Bh:  bh,
		Fsm: fsm_,
	}
}
