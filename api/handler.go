package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"kumparan/constants"
	"kumparan/repository"
	"kumparan/service"
)

type Handler interface {
	CreateNews(w http.ResponseWriter, r *http.Request)
	GetNews(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	producer service.ProducerService
	cache    repository.Cache
}

func InitHandler(producer service.ProducerService, cache repository.Cache) Handler {
	return &handler{producer, cache}
}

func (h *handler) CreateNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	err = h.producer.CreateNews(data)
	if err != nil {
		response := FailedResponse{
			Status:  http.StatusForbidden,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SuccessResponse{
		Status:  http.StatusOK,
		Message: "Successfully sent a message!",
	}
	json.NewEncoder(w).Encode(response)
}

func (h *handler) GetNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	content := h.cache.Get(r.RequestURI)
	if content != nil {
		log.Print("Cache Hit!\n")
		json.NewEncoder(w).Encode(SuccessResponse{
			Status:  http.StatusOK,
			Message: "Successfully retrieved news from cache!",
			Data:    content,
		})
		return
	}

	page, _ := strconv.Atoi(r.FormValue("page"))

	news, err := h.producer.GetAllNews(page)
	if err != nil {
		response := FailedResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SuccessResponse{
		Status:  http.StatusOK,
		Message: "Successfully retrieved news!",
		Data:    news,
	}

	if duration, err := time.ParseDuration(constants.CacheDuration); err == nil {
		log.Printf("New data cached: %s for %s\n", r.RequestURI, duration)
		h.cache.Set(r.RequestURI, news, duration)
	} else {
		log.Printf("Page not cached. err: %s\n", err)
	}

	json.NewEncoder(w).Encode(response)
}
