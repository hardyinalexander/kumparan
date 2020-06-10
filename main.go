package main

import (
	"fmt"
	"kumparan/api"
	"kumparan/config"
	"kumparan/repository"
	"kumparan/service"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	cfg := config.Get()

	db := initDB(cfg)
	defer db.Close()

	repo := repository.InitRepository(db)
	consumerSvc := service.InitConsumer(repo)
	router := initRouter(consumerSvc)

	fmt.Printf("Starting server on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router))
}

func initDB(cfg config.Config) *gorm.DB {
	connInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)
	db, err := gorm.Open("postgres", connInfo)
	handleError(err, "Failed to connect to database")

	return db
}

func initRouter(consumerSvc service.Consumer) *mux.Router {
	router := mux.NewRouter()
	handler := api.InitHandler(consumerSvc)

	router.HandleFunc("/news", handler.CreateNews).Methods("POST")

	return router
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}
