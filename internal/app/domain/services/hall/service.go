package hall

import (
	"context"
	"errors"
	"fmt"
	"movies/internal/app/core/microservices/booking"
	dto "movies/internal/app/domain/core/dto/hall"
	repository "movies/internal/app/domain/repositories/hall"
	"movies/internal/app/models"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("hall not found")
)

type Service interface {
	Create(ctx context.Context, input dto.CreateHallDTO) (models.Hall, error)
	Get(ctx context.Context, id int64) (models.Hall, error)
	Delete(ctx context.Context, id int64) (bool, error)
	List(ctx context.Context) ([]models.Hall, error)
}

type service struct {
	db            *gorm.DB
	repository    repository.Repository
	bookingClient booking.Client
}

func NewService(
	db *gorm.DB,
	bookingClient booking.Client,
) Service {
	return &service{
		db:            db,
		repository:    repository.NewRepository(db),
		bookingClient: bookingClient,
	}
}

func (s *service) Create(ctx context.Context, input dto.CreateHallDTO) (models.Hall, error) {
	hall, err := s.repository.Create(ctx, dto.CreateHallDTO{
		CinemaID: input.CinemaID,
		Name:     input.Name,
		Seats:    input.Seats,
	})
	if err != nil {
		return models.Hall{}, err
	}

	if err = s.bookingClient.StoreSeat(ctx, booking.StoreSeat{
		HallID: hall.ID,
		Row:    hall.Seats / 10,
		Number: hall.Seats / (hall.Seats / 10),
	}); err != nil {
		return models.Hall{}, fmt.Errorf("failed store seat: %w", err)
	}

	return hall, nil
}

func (s *service) Get(ctx context.Context, id int64) (models.Hall, error) {
	hall, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return hall, ErrNotFound
		}
	}

	return hall, err
}

func (s *service) Delete(ctx context.Context, id int64) (bool, error) {
	ok, err := s.repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ok, ErrNotFound
		}
	}

	if err = s.bookingClient.DeleteSeat(ctx, id); err != nil {
		return false, fmt.Errorf("failed delete seat: %w", err)
	}

	return ok, err
}

func (s *service) List(ctx context.Context) ([]models.Hall, error) {
	hall, err := s.repository.List(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return hall, ErrNotFound
		}
	}

	return hall, err
}
