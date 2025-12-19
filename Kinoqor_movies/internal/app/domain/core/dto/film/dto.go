package film

import "time"

type CreateFilmDTO struct {
	Name        string    `json:"name"  validate:"required"`
	Description string    `json:"description" binding:"omitempty"`
	StartDate   string    `json:"start_date" validate:"required"` 
	EndDate     string    `json:"end_date" validate:"required"`
	Duration    string    `json:"duration"  validate:"required"`
	Premier     string    `json:"premier"  validate:"required"`
	Production  string    `json:"production" validate:"required"`
	Director    string    `json:"director"  validate:"required"`
	Rate        float32   `json:"rate"  validate:"required"`
	AgeLimit    int8      `json:"age_limit"  validate:"required"`
}

type CreateFilmGenreDTO struct {
	FilmID        int64    `json:"film_id"  validate:"required"`
	GenreID       int64    `json:"genre_id"  validate:"required"`
}

type CreateFilmServiceDTO struct {
	Name        string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Duration    string
	Premier     time.Time
	Production  string
	Director    string
	Rate        float32
	AgeLimit    int8
}

type CreateFilmGenreServiceDTO struct {
	FilmID  int64
	GenreID int64 
}