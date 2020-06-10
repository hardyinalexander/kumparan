package api

import (
	"encoding/json"
	"net/http"

	"kumparan/service"
)

type Handler interface {
	CreateNews(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	consumer service.Consumer
}

func InitHandler(consumer service.Consumer) Handler {
	return &handler{consumer}
}

func (h *handler) CreateNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		response := FailedResponse{
			Status:  http.StatusForbidden,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	author := r.FormValue("author")
	body := r.FormValue("body")

	news, err := h.consumer.CreateNews(author, body)
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
		Message: "You have created a news",
		Data:    news,
	}
	json.NewEncoder(w).Encode(response)
}
