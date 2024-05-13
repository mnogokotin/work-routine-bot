package rabbitmq

import (
	"context"
	"encoding/gob"
	"github.com/mnogokotin/golang-packages/message_queue/rabbitmq"
	"github.com/mnogokotin/golang-packages/utils/e"
	ur "github.com/mnogokotin/golang-packages/utils/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"time"
	"work-routine-bot/internal/domain"
)

type Storage struct {
	*rabbitmq.Rabbitmq
}

func (s Storage) SendOnCreateMessage(ctx context.Context, task *domain.Task) error {
	queue, err := s.Ch.QueueDeclare("task-created", false, false, false, false, nil)
	if err != nil {
		return e.Wrap("failed to declare a task created queue", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gob.Register(task)
	body := make(map[string]interface{})
	body["task"] = task
	bodyEncoded := ur.Encode(body)

	err = s.Ch.PublishWithContext(ctx,
		"",
		queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(bodyEncoded),
		})
	if err != nil {
		return e.Wrap("failed to publish a task created message", err)
	}

	return nil
}
