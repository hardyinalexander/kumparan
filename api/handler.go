package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"kumparan/service"
)

type Handler interface {
	CreateNews(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	producer service.ProducerService
}

func InitHandler(producer service.ProducerService) Handler {
	return &handler{producer}
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
