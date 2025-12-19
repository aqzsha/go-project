package models

func (FilmGenre) TableName() string {
	return "film_genre"
}

type FilmGenre struct {
	FilmID  int64
	GenreID int64
}
