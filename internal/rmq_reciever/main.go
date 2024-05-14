package rmq_reciever

import (
	"encoding/gob"
	"github.com/mnogokotin/golang-packages/logger"
	"github.com/mnogokotin/golang-packages/message_queue/rabbitmq"
	d2 "github.com/mnogokotin/golang-packages/utils/d"
	ur "github.com/mnogokotin/golang-packages/utils/rabbitmq"
	"os"
	"os/signal"
	"syscall"
	"work-routine-bot/internal/config"
	"work-routine-bot/internal/domain"
)

func Run() {
	cfg := config.New()

	log := logger.New(cfg.Env)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{}, 1)

	rmq, err := rabbitmq.New(cfg.Rabbitmq.Uri)
	if err != nil {
		panic(err)
	}

	q, err := rmq.Ch.QueueDeclare(
		"task-created",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("failed to declare a reciever queue: " + err.Error())
	}

	msgs, err := rmq.Ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("failed to register a reciever: " + err.Error())
	}

	go func() {
		for d := range msgs {
			gob.Register(&domain.Task{})
			decodedBody := ur.Decode(string(d.Body))
			d2.Dnd(decodedBody["task"])

			log.Info("received a message: %s", decodedBody)
		}
	}()

	go func() {
		<-sigs

		rmq.Close()

		done <- struct{}{}
	}()

	log.Info("rmq reciever started")
	//bot_.Bh.Start()

	<-done
	log.Info("rmq reciever stopped")
}
