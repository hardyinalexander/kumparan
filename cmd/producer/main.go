package main

import (
	"fmt"
	"kumparan/api"
	"kumparan/config"
	"kumparan/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

func main() {
	cfg := config.Get()

	mqConn, err := amqp.Dial(cfg.MQURL)
	handleError(err, "Failed to connect to RabbitMQ")
	defer mqConn.Close()

	ch, err := mqConn.Channel()
	handleError(err, "Failed to open a channel")
	defer ch.Close()

	producerSvc := service.InitProducerService(ch)
	router := initRouter(producerSvc)

	fmt.Printf("Starting server on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router))
}

func initRouter(producerSvc service.ProducerService) *mux.Router {
	router := mux.NewRouter()
	handler := api.InitHandler(producerSvc)

	router.HandleFunc("/news", handler.CreateNews).Methods("POST")

	return router
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
