package seat

import (
	dto "booking/internal/app/domain/core/dto/seat"
	repository "booking/internal/app/domain/repositories/seat"
	"booking/internal/app/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("seat not found")
)

type Service interface {
	Create(ctx context.Context, input dto.CreateSeatDTO) (models.Seat, error)
	Delete(ctx context.Context, id int64) (bool, error)
}

type service struct {
	db         *gorm.DB
	repository repository.Repository
}

func NewService(
	db *gorm.DB,
) Service {
	return &service{
		db:         db,
		repository: repository.NewRepository(db),
	}
}

func (s *service) Create(ctx context.Context, payload dto.CreateSeatDTO) (models.Seat, error) {
	seat, err := s.repository.Create(ctx, dto.CreateSeatDTO{
		HallID: payload.HallID,
		Row:    payload.Row,
		Number: payload.Number,
	})
	if err != nil {
		return models.Seat{}, err
	}

	return seat, nil
}

func (s *service) Delete(ctx context.Context, id int64) (bool, error) {
	ok, err := s.repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ok, ErrNotFound
		}
	}

	return ok, err
}
