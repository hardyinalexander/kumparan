package service

import (
	"encoding/json"
	"kumparan/repository"
	"time"
)

type ConsumerService interface {
	CreateNews(data []byte) error
}

type consumerService struct {
	repo repository.Repository
	// elactic search
}

func InitConsumerService(repo repository.Repository) ConsumerService {
	return &consumerService{repo}
}

func (c *consumerService) CreateNews(data []byte) error {
	var n repository.News
	err := json.Unmarshal(data, &n)
	if err != nil {
		return err
	}

	n.Created = time.Now()
	err = c.repo.CreateNews(&n)
	return err
}
