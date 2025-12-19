package review

import (
	"context"
	"errors"
	"fmt"
	dto "movies/internal/app/domain/core/dto/film/review"
	"movies/internal/app/models"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("review not found")

type Repository interface {
	Create(ctx context.Context, input dto.ServiceCreateReviewDTO) (models.Review, error)
	Get(ctx context.Context, id int64) (models.Review, error)
	Delete(ctx context.Context, id int64) (bool, error)
	List(ctx context.Context) ([]models.Review, error)
	FilmList(ctx context.Context, id int64) ([]models.Review, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input dto.ServiceCreateReviewDTO) (models.Review, error) {
	review := models.Review{
		FilmID: input.FilmID,
		UserID: input.UserID,
		Body:   input.Body,
		Rating: input.Rating,
		Title:  input.Title,
	}

	if err := r.db.WithContext(ctx).Create(&review).Error; err != nil {
		return models.Review{}, fmt.Errorf("failed to create Review: %w", err)
	}

	return review, nil
}

func (r *repository) Get(ctx context.Context, id int64) (models.Review, error) {
	var review models.Review

	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&review, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return review, ErrNotFound
		}

		return review, fmt.Errorf("failed to get Review: %w", err)
	}

	return review, nil
}

func (r *repository) Delete(ctx context.Context, id int64) (bool, error) {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Review{})
	if res.Error != nil {
		return false, fmt.Errorf("failed to delete Review: %w", res.Error)
	}

	return true, nil
}

func (r *repository) List(ctx context.Context) ([]models.Review, error) {
	var review []models.Review

	if err := r.db.WithContext(ctx).
		Find(&review).Error; err != nil {

		return nil, fmt.Errorf("failed to get review list: %w", err)
	}

	return review, nil
}

func (r *repository) FilmList(ctx context.Context, id int64) ([]models.Review, error) {
	var review []models.Review

	if err := r.db.WithContext(ctx).Where("film_id = ?", id).Find(&review).Error; err != nil {
		return nil, fmt.Errorf("failed to get review list: %w", err)
	}

	return review, nil
}
