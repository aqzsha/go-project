package film

import "time"

type CreateFilmDTO struct {
	Name        string    `json:"name"  validate:"required"`
	Description string    `json:"description" binding:"omitempty"`
	DetailsID   int64     `json:"details_id"`
	StartDate   time.Time `json:"start_date" validate:"required"` 
	EndDate     time.Time `json:"end_date" validate:"required"`
	Duration    string    `json:"duration"  validate:"required"`
	Premier     time.Time `json:"premier"  validate:"required"`
	Production  string    `json:"production" validate:"required"`
	Director    string    `json:"director"  validate:"required"`
	Rate        float32   `json:"rate"  validate:"required"`
	AgeLimit    int8      `json:"age_limit"  validate:"required"`
}

type CreateFilmGenreDTO struct {
	FilmID        int64    `json:"film_id"  validate:"required"`
	GenreID       int64    `json:"genre_id"  validate:"required"`
}