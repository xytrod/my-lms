package broker

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	rabbit *RabbitMQ
}

func NewPublisher(rabbit *RabbitMQ) *Publisher {
	return &Publisher{rabbit: rabbit}
}
func (p *Publisher) Publish(ctx context.Context, routingKey string, event any) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.rabbit.ch.PublishWithContext(ctx, ExchangeName, routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		})
}
