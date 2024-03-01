package main

import (
	"log"
	ctelegram "work-routine-bot/internal/clients/telegram"
	"work-routine-bot/internal/config"
	"work-routine-bot/internal/consumer/event-consumer"
	"work-routine-bot/internal/events/telegram"
	"work-routine-bot/internal/storage/pages/mongo"
)

func main() {
	cfg := config.New()

	storage := mongo.New(cfg.Mongo.Uri, cfg.Mongo.ConnectTimeout, cfg.Mongo.DbName)

	eventsProcessor := telegram.New(
		ctelegram.New(cfg.Tg.Host, cfg.Tg.Token),
		storage,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, cfg.Bot.BatchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
