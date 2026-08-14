package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const ExchangeName = "lms.events"

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	err = ch.ExchangeDeclare(
		ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		_ = conn.Close()
		_ = ch.Close()
		return nil, err
	}
	return &RabbitMQ{conn: conn, ch: ch}, nil
}
func (r *RabbitMQ) Close() error {
	if r.ch != nil {
		_ = r.ch.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
