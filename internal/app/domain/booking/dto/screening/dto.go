package screening

type CreateScreeningDTO struct {
	CinemaID int64   `json:"cinema_id"  validate:"required"`
	HallID   int64   `json:"hall_id" validate:"required"`
	FilmID   int64   `json:"film_id"  validate:"required"`
	StartAt  string  `json:"start_at" validate:"required"`
	EndAt    string  `json:"end_at" validate:"required"`
	Language string  `json:"language"  validate:"required"`
	Format   string  `json:"format"  validate:"required"`
	Price    float64 `json:"price" validate:"required"`
}
