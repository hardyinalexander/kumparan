package service

import (
	"kumparan/repository"
	"time"
)

type Consumer interface {
	CreateNews(author, body string) (repository.News, error)
}

type consumer struct {
	repo repository.Repository
	// elactic search
}

func InitConsumer(repo repository.Repository) Consumer {
	return &consumer{repo}
}

func (c *consumer) CreateNews(author, body string) (repository.News, error) {
	n := repository.News{
		Author:  author,
		Body:    body,
		Created: time.Now(),
	}

	n, err := c.repo.CreateNews(n)

	return n, err
}
