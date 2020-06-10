package service

import (
	"log"

	"github.com/streadway/amqp"
)

type ProducerService interface {
	CreateNews(data []byte) error
}

type producerService struct {
	ch *amqp.Channel
}

func InitProducerService(ch *amqp.Channel) ProducerService {
	return &producerService{ch}
}

func (p *producerService) CreateNews(data []byte) error {
	q, err := p.ch.QueueDeclare(
		"news-queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	handleError(err, "Failed to declare a queue")

	err = p.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	handleError(err, "Failed to publish a message")

	return err
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
