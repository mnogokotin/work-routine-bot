package app

import (
	"github.com/mnogokotin/golang-packages/database/postgres"
	"github.com/mnogokotin/golang-packages/logger"
	"github.com/mnogokotin/golang-packages/message_queue/rabbitmq"
	"os"
	"os/signal"
	"syscall"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/config"
	apptg "work-routine-bot/internal/handler/app/tg"
	ttg "work-routine-bot/internal/handler/task/tg"
	ppg "work-routine-bot/internal/storage/projects/postgres"
	tpg "work-routine-bot/internal/storage/tasks/postgres"
	trmq "work-routine-bot/internal/storage/tasks/rabbitmq"
	upg "work-routine-bot/internal/storage/users/postgres"
)

func Run() {
	cfg := config.New()

	log := logger.New(cfg.Env)

	bot_ := bot.New(cfg.Tg.Token, cfg.Env)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{}, 1)

	pg, err := postgres.New(cfg.Postgres.Uri)
	if err != nil {
		panic(err)
	}

	rmq, err := rabbitmq.New(cfg.Rabbitmq.Uri)
	if err != nil {
		panic(err)
	}

	projectPqStorage := &ppg.Storage{
		Postgres: pg,
	}
	taskPgStorage := &tpg.Storage{
		Postgres: pg,
	}
	taskRmqStorage := &trmq.Storage{
		Rabbitmq: rmq,
	}
	userPgStorage := &upg.Storage{
		Postgres: pg,
	}

	appHandler := apptg.New(log, bot_)
	appHandler.Handle()
	defer appHandler.HandleEnd()

	ttg.New(log, bot_, projectPqStorage, taskPgStorage, taskRmqStorage, userPgStorage).Handle()

	go func() {
		<-sigs

		bot_.Bot.StopLongPolling()
		bot_.Bh.Stop()
		pg.Close()
		rmq.Close()

		done <- struct{}{}
	}()

	log.Info("bot started")
	bot_.Bh.Start()

	<-done
	log.Info("bot stopped")
}
