package main

import (
	"fmt"
	"kumparan/api"
	"kumparan/config"
	"kumparan/repository"
	"kumparan/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
	elastic "gopkg.in/olivere/elastic.v6"
)

func main() {
	cfg := config.Get()

	mqConn, err := amqp.Dial(cfg.MQURL)
	handleError(err, "Failed to connect to RabbitMQ")
	defer mqConn.Close()

	ch, err := mqConn.Channel()
	handleError(err, "Failed to open a channel")
	defer ch.Close()

	db := initDB(cfg)
	defer db.Close()

	e := initElasticSearch(cfg)

	repo := repository.InitRepository(db, e)

	producerSvc := service.InitProducerService(ch, repo)
	router := initRouter(producerSvc)

	log.Printf("Starting server on port %s... \n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router))
}

func initRouter(producerSvc service.ProducerService) *mux.Router {
	router := mux.NewRouter()
	handler := api.InitHandler(producerSvc)

	router.HandleFunc("/news", handler.CreateNews).Methods("POST")
	router.HandleFunc("/news", handler.GetNews).Queries("page", "{page:[0-9]+}").Methods("GET")

	return router
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
	return client
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
