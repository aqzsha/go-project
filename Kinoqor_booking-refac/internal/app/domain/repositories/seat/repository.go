package seat

import (
	dto "booking/internal/app/domain/core/dto/seat"
	"booking/internal/app/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("seat not found")

type Repository interface {
	Create(ctx context.Context, input dto.CreateSeatDTO) (models.Seat, error)
	Delete(ctx context.Context, id int64) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input dto.CreateSeatDTO) (models.Seat, error) {
	seat := models.Seat{
		HallID: input.HallID,
		Row:    input.Row,
		Number: input.Number,
	}

	if err := r.db.WithContext(ctx).Create(&seat).Error; err != nil {
		return models.Seat{}, fmt.Errorf("failed to create Seat: %w", err)
	}

	return seat, nil
}

func (r *repository) Delete(ctx context.Context, id int64) (bool, error) {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Seat{})
	if res.Error != nil {
		return false, fmt.Errorf("failed to delete Seat: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return false, ErrNotFound
	}

	return true, nil
}
