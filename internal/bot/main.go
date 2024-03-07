package bot

import (
	"github.com/mnogokotin/golang-packages/logger"
	"github.com/mymmrac/telego"
	"work-routine-bot/internal/config"
	"work-routine-bot/internal/consumer/update-consumer"
	"work-routine-bot/internal/processor/pages/tg"
	"work-routine-bot/internal/storage/pages/mongo"
)

func Run() {
	cfg := config.New()

	log := logger.New(cfg.Env)

	pagesStorage := mongo.New(cfg.Mongo.Uri, cfg.Mongo.ConnectTimeout, cfg.Mongo.DbName)

	bot, err := telego.NewBot(cfg.Tg.Token)
	if err != nil {
		panic("can't create bot: " + err.Error())
	}

	updatesChan, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		panic("can't create bot's updates channel: " + err.Error())
	}
	defer bot.StopLongPolling()

	pagesProcessor := tg.New(log, bot, updatesChan, pagesStorage)

	log.Info("service started")

	consumer := update_consumer.New(log, pagesProcessor, pagesProcessor)

	if err = consumer.Start(); err != nil {
		log.Error("service is stopped", err)
	}
}
