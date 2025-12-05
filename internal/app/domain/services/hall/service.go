package hall

import (
	"context"
	"errors"
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

func (s *service) Create(ctx context.Context, input dto.CreateHallDTO) (models.Hall, error) {
	hall, err := s.repository.Create(ctx, dto.CreateHallDTO{
		CinemaID: input.CinemaID,
		Name:     input.Name,
		Seats:    input.Seats,
	})
	if err != nil {
		return models.Hall{}, err
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