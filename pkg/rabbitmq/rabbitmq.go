package rabbitmq

import (
	"fmt"
)

type Rabbitmq struct {
	Conn *amqp091.Connection
	Ch   *amqp091.Channel
}

func New(connectionUri string) (*Rabbitmq, error) {
	conn, err := amqp091.Dial(connectionUri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Rabbitmq{
		Conn: conn,
		Ch:   ch,
	}, nil
}

func (p *Rabbitmq) Close() {
	if err := p.Ch.Close(); err != nil {
		fmt.Printf("could not close rabbinmq channel: %s", err)
	}
	if err := p.Conn.Close(); err != nil {
		fmt.Printf("could not close rabbinmq connection: %s", err)
	}
}
