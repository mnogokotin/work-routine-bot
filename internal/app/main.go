package app

import (
	"github.com/mnogokotin/golang-packages/database/postgres"
	"github.com/mnogokotin/golang-packages/logger"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/config"
	apptg "work-routine-bot/internal/handler/app/tg"
	whtg "work-routine-bot/internal/handler/working-hours/tg"
	tpostgres "work-routine-bot/internal/storage/tasks/postgres"
	upostgres "work-routine-bot/internal/storage/users/postgres"
)

func Run() {
	cfg := config.New()

	log := logger.New(cfg.Env)

	bot_ := bot.New(cfg.Tg.Token)
	defer bot_.Bh.Stop()
	defer bot_.Bot.StopLongPolling()

	log.Info("service started")

	ppg, err := postgres.New(cfg.Postgres.Uri)
	if err != nil {
		panic(err)
	}

	userStorage := &upostgres.Storage{
		Postgres: ppg,
	}
	taskStorage := &tpostgres.Storage{
		Postgres: ppg,
	}

	apptg.New(log, bot_).Handle()
	whtg.New(log, bot_, userStorage, taskStorage).Handle()

	bot_.Bh.Start()
}
