package film

type CreateFilmDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	GenreID     int64  `json:"genre_id"`
}

type FilmIDParam struct {
	ID int64 `uri:"id"`
}
