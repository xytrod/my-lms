package broker

import (
	"encoding/json"
	"log"
)

const UserEnrolledQueue = "course.user-enrolled"

func (r *RabbitMQ) StartUserEnrolledConsumer() error {
	queue, err := r.ch.QueueDeclare(UserEnrolledQueue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = r.ch.QueueBind(queue.Name, "user.enrolled", ExchangeName, false, nil)
	if err != nil {
		return err
	}
	messages, err := r.ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for message := range messages {
			var event UserEnrolledEvent
			if err := json.Unmarshal(message.Body, &event); err != nil {
				log.Printf("invalid user-enrolled event: %s\n", err)
				if err := message.Nack(false, false); err != nil {
					log.Printf("failed to nack user-enrolled event: %s\n", err)
				}

				continue
			}
			log.Printf("user-enrolled event: user_id=%s course_id=%s\n", event.UserID, event.CourseID)
			if err := message.Ack(false); err != nil {
				log.Printf("failed to ack user-enrolled event: %s\n", err)
			}
		}
	}()
	return nil
}
