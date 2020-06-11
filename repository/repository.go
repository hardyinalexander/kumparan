package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"kumparan/constants"

	"github.com/jinzhu/gorm"
	elastic "gopkg.in/olivere/elastic.v6"
)

type Repository interface {
	CreateNews(news *News) error
	CreateESDocument(id int, created string) error
	GetNewsByID(id int) (*News, error)
	GetAllESDocument(page int) ([]int, error)
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

func (r *repository) GetNewsByID(id int) (*News, error) {
	news := &News{}
	err := r.db.Where("id = ?", id).First(&news).Error

	return news, err
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

func (r *repository) GetAllESDocument(page int) ([]int, error) {
	ctx := context.Background()
	from := (page - 1) * constants.PageSize
	var result []int
	query := elastic.MatchAllQuery{}

	searchResult, err := r.e.Search().
		TrackTotalHits(false).
		Index(constants.IndexName).
		Query(query).
		Sort("created", false).
		From(from).
		Size(constants.PageSize).
		Pretty(true).
		Do(ctx)
	if err != nil {
		return result, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var n News
		err = json.Unmarshal(*hit.Source, &n)
		if err != nil {
			return result, err
		}
		result = append(result, n.ID)
	}
	return result, err
}
