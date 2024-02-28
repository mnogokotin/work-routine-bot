package main

import (
	"log"
	"time"
	tgClient "work-routine-bot/internal/clients/telegram"
	"work-routine-bot/internal/config"
	"work-routine-bot/internal/consumer/event-consumer"
	"work-routine-bot/internal/events/telegram"
	"work-routine-bot/internal/storage/mongo"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	cfg := config.MustLoad()

	storage := mongo.New(cfg.MongoConnectionString, 10*time.Second)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, cfg.TgBotToken),
		storage,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
