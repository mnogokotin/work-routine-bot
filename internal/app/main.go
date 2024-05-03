package app

import (
	"github.com/mnogokotin/golang-packages/database/postgres"
	"github.com/mnogokotin/golang-packages/logger"
	"os"
	"os/signal"
	"syscall"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/config"
	apptg "work-routine-bot/internal/handler/app/tg"
	ttg "work-routine-bot/internal/handler/task/tg"
	ppostgres "work-routine-bot/internal/storage/projects/postgres"
	tpostgres "work-routine-bot/internal/storage/tasks/postgres"
	upostgres "work-routine-bot/internal/storage/users/postgres"
)

func Run() {
	cfg := config.New()

	log := logger.New(cfg.Env)

	bot_ := bot.New(cfg.Tg.Token, cfg.Env)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{}, 1)

	ppg, err := postgres.New(cfg.Postgres.Uri)
	if err != nil {
		panic(err)
	}

	projectStorage := &ppostgres.Storage{
		Postgres: ppg,
	}
	taskStorage := &tpostgres.Storage{
		Postgres: ppg,
	}
	userStorage := &upostgres.Storage{
		Postgres: ppg,
	}

	appHandler := apptg.New(log, bot_)
	appHandler.Handle()
	defer appHandler.HandleEnd()

	ttg.New(log, bot_, projectStorage, taskStorage, userStorage).Handle()

	go func() {
		<-sigs

		bot_.Bot.StopLongPolling()
		bot_.Bh.Stop()

		done <- struct{}{}
	}()

	log.Info("bot started")
	bot_.Bh.Start()

	<-done
	log.Info("bot stopped")
}
