package genre

import (
	"context"
	"errors"
	"fmt"
	dto "movies/internal/app/domain/core/dto/genre"
	"movies/internal/app/models"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("genre not found")

type Repository interface {
	Create(ctx context.Context, input dto.CreateGenreDTO) (models.Genre, error)
	Get(ctx context.Context, id int64) (models.Genre, error)
	Delete(ctx context.Context, id int64) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input dto.CreateGenreDTO) (models.Genre, error) {
	genre := models.Genre{
		Name: input.Name,
	}

	if err := r.db.WithContext(ctx).Create(&genre).Error; err != nil {
		return models.Genre{}, fmt.Errorf("failed to create Genre: %w", err)
	}


	return genre, nil
}

func (r *repository) Get(ctx context.Context, id int64) (models.Genre, error) {
	var genre models.Genre

	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&genre, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return genre, ErrNotFound
		}

		return genre, fmt.Errorf("failed to get Genre: %w", err)
	}

	return genre, nil
}

func (r *repository) Delete(ctx context.Context, id int64) (bool, error) {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Genre{})
	if res.Error != nil {
		return false, fmt.Errorf("failed to delete Genre: %w", res.Error)
	}

	return true, nil
}


