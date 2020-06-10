package repository

import (
	"time"
)

type News struct {
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}
