package app

import (
	ppostgres "github.com/mnogokotin/golang-packages/database/postgres"
	"github.com/mnogokotin/golang-packages/logger"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/config"
	"work-routine-bot/internal/consumer/update-consumer"
	"work-routine-bot/internal/processor/working-hours/tg"
	tpostgres "work-routine-bot/internal/storage/tasks/postgres"
	upostgres "work-routine-bot/internal/storage/users/postgres"
)

func Run() {
	cfg := config.New()

	log := logger.New(cfg.Env)

	bot_ := bot.New(cfg.Tg.Token)
	defer bot_.Bot.StopLongPolling()

	ppg, err := ppostgres.New(cfg.Postgres.Uri)
	if err != nil {
		panic(err)
	}
	userStorage := &upostgres.Storage{
		Postgres: ppg,
	}
	taskStorage := &tpostgres.Storage{
		Postgres: ppg,
	}

	workingHoursProcessor := tg.New(log, bot_, userStorage, taskStorage)

	log.Info("service started")

	consumer := update_consumer.New(log, workingHoursProcessor, workingHoursProcessor)

	if err := consumer.Start(); err != nil {
		log.Error("service is stopped", err)
	}
}
