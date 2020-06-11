package repository

import (
	"context"
	"fmt"
	"kumparan/constants"

	"github.com/jinzhu/gorm"
	elastic "gopkg.in/olivere/elastic.v6"
)

type Repository interface {
	CreateNews(news *News) error
	CreateESDocument(id int, created string) error
}

type repository struct {
	db *gorm.DB
	e  *elastic.Client
}

func InitRepository(db *gorm.DB, e *elastic.Client) Repository {
	return &repository{db, e}
}

func (r *repository) CreateNews(news *News) error {
	err := r.db.Create(&news).Error
	return err
}

func (r *repository) CreateESDocument(id int, created string) error {
	ctx := context.Background()
	body := fmt.Sprintf(`{
		"id": %d,
		"created": "%s"
	}`, id, created)

	_, err := r.e.Index().
		Index(constants.IndexName).
		Type(constants.DocType).
		BodyJson(body).
		Do(ctx)

	if err != nil {
		return err
	}

	// Flush to make sure the documents got written.
	_, err = r.e.Flush().Index(constants.IndexName).Do(ctx)
	return err
}
