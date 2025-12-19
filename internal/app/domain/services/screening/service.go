package screening

import (
	dto "booking/internal/app/domain/core/dto/screening"
	repository "booking/internal/app/domain/repositories/screening"
	"booking/internal/app/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("screening not found")
)

type Service interface {
	Create(ctx context.Context, input dto.ServiceCreateScreeningDTO) (models.Screening, error)
	Get(ctx context.Context, id int64) (models.Screening, error)
	Delete(ctx context.Context, id int64) (bool, error)
	ListByFilm(ctx context.Context, filmId int64) ([]models.Screening, error)
	ListByCinema(ctx context.Context, cinemaId int64) ([]models.Screening, error)
	ListByHall(ctx context.Context, hallId int64) ([]models.Screening, error)
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

func (s *service) Create(ctx context.Context, payload dto.ServiceCreateScreeningDTO) (models.Screening, error) {
	screening, err := s.repository.Create(ctx, dto.ServiceCreateScreeningDTO{
		CinemaID: payload.CinemaID,
		HallID:   payload.HallID,
		FilmID:   payload.FilmID,
		Price:    payload.Price,
		Format:   payload.Format,
		Language: payload.Language,
		StartAt:  payload.StartAt,
		EndAt:    payload.EndAt,
	})
	if err != nil {
		return models.Screening{}, err
	}

	return screening, nil
}

func (s *service) Get(ctx context.Context, id int64) (models.Screening, error) {
	screening, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return screening, ErrNotFound
		}
	}

	return screening, err
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

func (s *service) ListByFilm(ctx context.Context, filmId int64) ([]models.Screening, error) {
	screening, err := s.repository.ListByFilm(ctx, filmId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return screening, ErrNotFound
		}
	}

	return screening, err
}

func (s *service) ListByCinema(ctx context.Context, cinemaId int64) ([]models.Screening, error) {
	screening, err := s.repository.ListByCinema(ctx, cinemaId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return screening, ErrNotFound
		}
	}

	return screening, err
}

func (s *service) ListByHall(ctx context.Context, hallId int64) ([]models.Screening, error) {
	screening, err := s.repository.ListByHall(ctx, hallId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return screening, ErrNotFound
		}
	}

	return screening, err
}
