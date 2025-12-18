package cinema

import (
	"context"
	"errors"
	dto "movies/internal/app/domain/core/dto/cinema"
	repository "movies/internal/app/domain/repositories/cinema"
	"movies/internal/app/models"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("cinema not found")
)

type Service interface {
	Create(ctx context.Context, input dto.CreateCinemaDTO) (models.Cinema, error)
	Get(ctx context.Context, id int64) (models.Cinema, error)
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

func (s *service) Create(ctx context.Context, input dto.CreateCinemaDTO) (models.Cinema, error) {
	cinema, err := s.repository.Create(ctx, dto.CreateCinemaDTO{
		Name:        input.Name,
		Description: input.Description,
		Address: input.Address,
		Longitude: input.Longitude,
		Latitude: input.Latitude,
	})
	if err != nil {
		return models.Cinema{}, err
	}

	return cinema, nil
}

func (s *service) Get(ctx context.Context, id int64) (models.Cinema, error) {
	cinema, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return cinema, ErrNotFound
		}
	}

	return cinema, err
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
