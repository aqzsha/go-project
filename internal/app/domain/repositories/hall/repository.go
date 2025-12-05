package hall

import (
	"context"
	"errors"
	"fmt"
	dto "movies/internal/app/domain/core/dto/hall"
	"movies/internal/app/models"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("hall not found")

type Repository interface {
	Create(ctx context.Context, input dto.CreateHallDTO) (models.Hall, error)
	Get(ctx context.Context, id int64) (models.Hall, error)
	Delete(ctx context.Context, id int64) (bool, error)
	List(ctx context.Context) ([]models.Hall, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input dto.CreateHallDTO) (models.Hall, error) {
	hall := models.Hall{
		CinemaID: input.CinemaID,
		Name:     input.Name,
		Seats:    input.Seats,
	}

	if err := r.db.WithContext(ctx).Create(&hall).Error; err != nil {
		return models.Hall{}, fmt.Errorf("failed to create Hall: %w", err)
	}


	return hall, nil
}

func (r *repository) Get(ctx context.Context, id int64) (models.Hall, error) {
	var hall models.Hall

	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&hall, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return hall, ErrNotFound
		}

		return hall, fmt.Errorf("failed to get Hall: %w", err)
	}

	return hall, nil
}

func (r *repository) Delete(ctx context.Context, id int64) (bool, error) {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Hall{})
	if res.Error != nil {
		return false, fmt.Errorf("failed to delete Hall: %w", res.Error)
	}

	return true, nil
}

func (r *repository) List(ctx context.Context) ([]models.Hall, error) {
	var hall []models.Hall

	if err := r.db.WithContext(ctx).
		Find(&hall).Error; err != nil {

		return nil, fmt.Errorf("failed to get hall list: %w", err)
	}

	return hall, nil
}