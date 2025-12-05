package film

import (
	"context"
	"errors"
	"fmt"
	dto "movies/internal/app/domain/core/dto/film"
	"movies/internal/app/models"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("film not found")

type Repository interface {
	Create(ctx context.Context, input dto.CreateFilmDTO) (models.Film, error)
	Get(ctx context.Context, id int64) (models.Film, error)
	Delete(ctx context.Context, id int64) (bool, error)
	CreateGenre(ctx context.Context, input dto.CreateFilmGenreDTO) (models.FilmGenre, error)
	DeleteGenre(ctx context.Context, id int64) (bool, error)
	GetGenreList(ctx context.Context) ([]models.FilmGenre, error)
 	List(ctx context.Context) ([]models.Film, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input dto.CreateFilmDTO) (models.Film, error) {
	filmDetails := models.FilmDetails{
		Duration: input.Duration,
		Premier: input.Premier,
		Production: input.Production,
		Director: input.Director,
		Rate: float64(input.Rate),
		AgeLimit: int(input.AgeLimit),
	}

	if err := r.db.WithContext(ctx).Create(&filmDetails).Error; err != nil {
		return models.Film{}, fmt.Errorf("failed to create Film Details: %w", err)
	}

	film := models.Film{
		Name: input.Name,
		Description: input.Description,
		DetailsID: filmDetails.ID,
		StartDate: input.StartDate,
		EndDate: input.EndDate,
	}

	if err := r.db.WithContext(ctx).Create(&film).Error; err != nil {
		return models.Film{}, fmt.Errorf("failed to create Film: %w", err)
	}


	return film, nil
}

func (r *repository) Get(ctx context.Context, id int64) (models.Film, error) {
	var film models.Film

	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&film, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return film, ErrNotFound
		}

		return film, fmt.Errorf("failed to get Film: %w", err)
	}

	return film, nil
}

func (r *repository) Delete(ctx context.Context, id int64) (bool, error) {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Film{})
	if res.Error != nil {
		return false, fmt.Errorf("failed to delete Film: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return false, ErrNotFound
	}

	return true, nil
}


func (r *repository) CreateGenre(ctx context.Context, input dto.CreateFilmGenreDTO) (models.FilmGenre, error) {
	filmGenre := models.FilmGenre{
		FilmID: input.FilmID,
		GenreID: input.GenreID,
	}

	if err := r.db.WithContext(ctx).Create(&filmGenre).Error; err != nil {
		return models.FilmGenre{}, fmt.Errorf("failed to create Film Genre: %w", err)
	}

	return filmGenre, nil
}

	
func (r *repository) DeleteGenre(ctx context.Context, id int64) (bool, error) {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.FilmGenre{})
	if res.Error != nil {
		return false, fmt.Errorf("failed to delete Film Genre: %w", res.Error)
	}

	return true, nil
}

func (r *repository) GetGenreList(ctx context.Context) ([]models.FilmGenre, error) {
	var genres []models.FilmGenre

	if err := r.db.WithContext(ctx).
		Find(&genres).Error; err != nil {

		return nil, fmt.Errorf("failed to get genre list: %w", err)
	}

	return genres, nil
}

func (r *repository) List(ctx context.Context) ([]models.Film, error) {
	var film []models.Film

	if err := r.db.WithContext(ctx).
		Find(&film).Error; err != nil {

		return nil, fmt.Errorf("failed to get Film list: %w", err)
	}

	return film, nil
}