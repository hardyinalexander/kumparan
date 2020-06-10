package repository

import "github.com/jinzhu/gorm"

type Repository interface {
	CreateNews(news *News) error
}

type repository struct {
	db *gorm.DB
}

func InitRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateNews(news *News) error {
	err := r.db.Create(&news).Error
	return err
}
