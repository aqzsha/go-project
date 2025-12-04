package genre

import (
	"context"
	"errors"
	dto "movies/internal/app/domain/core/dto/genre"
	repository "movies/internal/app/domain/repositories/genre"
	"movies/internal/app/models"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("genre not found")
)

type Service interface {
	Create(ctx context.Context, input dto.CreateGenreDTO) (models.Genre, error)
	Get(ctx context.Context, id int64) (models.Genre, error)
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

func (s *service) Create(ctx context.Context, input dto.CreateGenreDTO) (models.Genre, error) {
	genre, err := s.repository.Create(ctx, dto.CreateGenreDTO{
		Name:        input.Name,
	})
	if err != nil {
		return models.Genre{}, err
	}

	return genre, nil
}

func (s *service) Get(ctx context.Context, id int64) (models.Genre, error) {
	genre, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return genre, ErrNotFound
		}
	}

	return genre, err
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
