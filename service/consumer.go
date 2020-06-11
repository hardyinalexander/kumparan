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

	n.Created = time.Now().Format("2006-01-02T15:04:05")
	err = c.repo.CreateNews(&n)
	if err != nil {
		return err
	}

	err = c.repo.CreateESDocument(n.ID, n.Created)
	return err
}
