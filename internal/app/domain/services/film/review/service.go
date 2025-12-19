package review

import (
	"context"
	"errors"
	dto "movies/internal/app/domain/core/dto/film/review"
	repository "movies/internal/app/domain/repositories/film/review"
	"movies/internal/app/models"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("film review not found")
)

type Service interface {
	Create(ctx context.Context, input dto.ServiceCreateReviewDTO) (models.Review, error)
	Get(ctx context.Context, id int64) (models.Review, error)
	Delete(ctx context.Context, id int64) (bool, error)
	List(ctx context.Context) ([]models.Review, error)
	FilmList(ctx context.Context, filmId int64) ([]models.Review, error)
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

func (s *service) Create(ctx context.Context, input dto.ServiceCreateReviewDTO) (models.Review, error) {
	review, err := s.repository.Create(ctx, dto.ServiceCreateReviewDTO{
		FilmID: input.FilmID,
		UserID: input.UserID,
		Body:   input.Body,
		Title:  input.Title,
		Rating: input.Rating,
	})
	if err != nil {
		return models.Review{}, err
	}

	return review, nil
}

func (s *service) Get(ctx context.Context, id int64) (models.Review, error) {
	review, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return review, ErrNotFound
		}
	}

	return review, err
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

func (s *service) List(ctx context.Context) ([]models.Review, error) {
	review, err := s.repository.List(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return review, ErrNotFound
		}
	}

	return review, err
}

func (s *service) FilmList(ctx context.Context, filmId int64) ([]models.Review, error) {
	review, err := s.repository.FilmList(ctx, filmId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return review, ErrNotFound
		}
	}

	return review, err
}
