package film

import (
	"context"
	"errors"
	dto "movies/internal/app/domain/core/dto/film"
	repository "movies/internal/app/domain/repositories/film"
	"movies/internal/app/models"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("movie not found")
	ErrGenreNotFound = errors.New("film genre not found")
)

type Service interface {
	Create(ctx context.Context, input dto.CreateFilmServiceDTO) (models.Film, error)
	Get(ctx context.Context, id int64) (models.Film, error)
	Delete(ctx context.Context, id int64) (bool, error)
	List(ctx context.Context) ([]models.Film, error)
	CreateGenre(ctx context.Context, input dto.CreateFilmGenreDTO) (models.FilmGenre, error)
	DeleteGenre(ctx context.Context, id int64) (bool, error)
	GetGenreList(ctx context.Context) ([]models.FilmGenre, error)
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

func (s *service) Create(ctx context.Context, input dto.CreateFilmServiceDTO) (models.Film, error) {
	film, err := s.repository.Create(ctx, dto.CreateFilmServiceDTO{
		Name:        input.Name,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Duration:    input.Duration,
		Premier:     input.Premier,
		Production:  input.Production,
		Director:    input.Director,
		Rate:        input.Rate,
		AgeLimit:    input.AgeLimit,
	})
	if err != nil {
		return models.Film{}, err
	}

	return film, nil
}

func (s *service) Get(ctx context.Context, id int64) (models.Film, error) {
	film, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return film, ErrNotFound
		}
	}

	return film, err
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

func (s *service) CreateGenre(ctx context.Context, input dto.CreateFilmGenreDTO) (models.FilmGenre, error) {
	genre, err := s.repository.CreateGenre(ctx, dto.CreateFilmGenreDTO{
		FilmID: input.FilmID,
		GenreID: input.GenreID,
	})
	if err != nil {
		return models.FilmGenre{}, err
	}

	return genre, nil
}

func (s *service) DeleteGenre(ctx context.Context, id int64) (bool, error) {
	ok, err := s.repository.DeleteGenre(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ok, ErrGenreNotFound
		}
	}

	return ok, err
}


func (s *service) GetGenreList(ctx context.Context) ([]models.FilmGenre, error) {
	film, err := s.repository.GetGenreList(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return film, ErrGenreNotFound
		}
	}

	return film, err
}


func (s *service) List(ctx context.Context) ([]models.Film, error) {
	review, err := s.repository.List(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return review, ErrNotFound
		}
	}

	return review, err
}