package film

import (
	"context"
	"errors"
	"fmt"
	dto "movies/internal/app/domain/core/dto/cinema"
	"movies/internal/app/models"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("film not found")

type Repository interface {
	Create(ctx context.Context, input dto.CreateCinemaDTO) (models.Cinema, error)
	Get(ctx context.Context, id int64) (models.Cinema, error)
	Delete(ctx context.Context, id int64) (bool, error)
	List(ctx context.Context) ([]models.Cinema, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input dto.CreateCinemaDTO) (models.Cinema, error) {
	cinemaDetails := models.CinemaDetails{
		Name: input.Name,
		Description: input.Description,
		Address: input.Address,
		Latitude: input.Latitude,
		Longitude: input.Longitude,
	}

	if err := r.db.WithContext(ctx).Create(&cinemaDetails).Error; err != nil {
		return models.Cinema{}, fmt.Errorf("failed to create Cinema Details: %w", err)
	}

	cinema := models.Cinema{
		Name: input.Name,
		DetailsID: cinemaDetails.ID,
		Details: cinemaDetails,
	}


	if err := r.db.WithContext(ctx).Create(&cinema).Error; err != nil {
		return models.Cinema{}, fmt.Errorf("failed to create Cinema: %w", err)
	}


	return cinema, nil
}

func (r *repository) Get(ctx context.Context, id int64) (models.Cinema, error) {
	var cinema models.Cinema

	if err := r.db.WithContext(ctx).
		Preload("Details").
		First(&cinema, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cinema, ErrNotFound
		}
		return cinema, fmt.Errorf("failed to get Cinema: %w", err)
	}

	return cinema, nil
}

func (r *repository) Delete(ctx context.Context, id int64) (bool, error) {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Cinema{})
	if res.Error != nil {
		return false, fmt.Errorf("failed to delete Cinema: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return false, ErrNotFound
	}

	return true, nil
}

func (r *repository) List(ctx context.Context) ([]models.Cinema, error) {
	var films []models.Cinema

	if err := r.db.WithContext(ctx).
		Preload("Details").
		Find(&films).Error; err != nil {
		return nil, fmt.Errorf("failed to get Cinema list: %w", err)
	}

	return films, nil
}