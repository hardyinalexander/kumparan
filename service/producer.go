package service

import (
	"kumparan/repository"
	"log"
	"sort"
	"sync"

	"github.com/streadway/amqp"
)

type ProducerService interface {
	CreateNews(data []byte) error
	GetAllNews(page int) ([]*repository.News, error)
}

type producerService struct {
	ch   *amqp.Channel
	repo repository.Repository
}

func InitProducerService(ch *amqp.Channel, repo repository.Repository) ProducerService {
	return &producerService{ch, repo}
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

func (p *producerService) GetAllNews(page int) ([]*repository.News, error) {
	ids, err := p.repo.GetAllESDocument(page)
	handleError(err, "Failed to retrieve documents from ElasticSearch")

	result := make([]*repository.News, len(ids))

	var wg sync.WaitGroup
	wg.Add(len(ids))

	for pos, id := range ids {
		go func(pos, id int) {
			defer wg.Done()
			news, err := p.repo.GetNewsByID(id)
			handleError(err, "Failed to GetNewsByID")

			result[pos] = news
		}(pos, id)
	}

	wg.Wait()
	sort.Sort(repository.NewsSorter(result))
	return result, err
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
