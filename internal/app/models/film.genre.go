package models

func (FilmGenre) TableName() string {
	return "film_genre"
}

type FilmGenre struct {
	FilmID  int64 `gorm:"column:film_id;primaryKey" json:"film_id"`
	GenreID int64 `gorm:"column:genre_id;primaryKey" json:"genre_id"`
}
