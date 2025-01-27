package main

import (
	"context"
	"errors"
	"fmt"
	"kumparan/config"
	"kumparan/constants"
	"kumparan/repository"
	"kumparan/service"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
	elastic "gopkg.in/olivere/elastic.v6"
)

func main() {
	cfg := config.Get()

	db := initDB(cfg)
	defer db.Close()

	e := initElasticSearch(cfg)

	repo := repository.InitRepository(db, e)
	consumerSvc := service.InitConsumerService(repo)

	conn, err := amqp.Dial(cfg.MQURL)
	handleError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"news-queue", // name
		false,        // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	handleError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	handleError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err = consumerSvc.CreateNews(d.Body)
			handleError(err, "Failed to create news from a message")
		}
	}()

	log.Printf("Waiting for messages...")
	<-forever

}

func initDB(cfg config.Config) *gorm.DB {
	connInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)
	db, err := gorm.Open("postgres", connInfo)
	handleError(err, "Failed to connect to Postgres")

	return db
}

func initElasticSearch(cfg config.Config) *elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL(cfg.ElasticSearchURL))

	handleError(err, "Failed to connect to ElasticSearch")
	go createIndexIfNotExist(client)
	return client
}

func createIndexIfNotExist(client *elastic.Client) {
	ctx := context.Background()

	exists, err := client.IndexExists(constants.IndexName).Do(ctx)
	handleError(err, "Failed checking if index exists")

	if exists {
		return
	}

	res, err := client.CreateIndex(constants.IndexName).Body(constants.IndexMapping).Do(ctx)
	handleError(err, "Failed creating index")

	if !res.Acknowledged {
		err = errors.New("CreateIndex was not acknowledged. Check that timeout value is correct.")
		handleError(err, err.Error())
	}
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
