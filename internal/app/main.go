package app

import (
	"github.com/mnogokotin/golang-packages/logger"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/config"
	"work-routine-bot/internal/consumer/update-consumer"
	"work-routine-bot/internal/processor/pages/tg"
	"work-routine-bot/internal/storage/pages/mongo"
)

func Run() {
	cfg := config.New()

	log := logger.New(cfg.Env)

	bot_ := bot.New(cfg.Tg.Token)
	defer bot_.Bot.StopLongPolling()

	pagesStorage := mongo.New(cfg.Mongo.Uri, cfg.Mongo.ConnectTimeout, cfg.Mongo.DbName)

	pagesProcessor := tg.New(log, bot_, pagesStorage)

	log.Info("service started")

	consumer := update_consumer.New(log, pagesProcessor, pagesProcessor)

	if err := consumer.Start(); err != nil {
		log.Error("service is stopped", err)
	}
}
