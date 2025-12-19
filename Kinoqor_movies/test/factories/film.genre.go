package factories

import (
	"movies/internal/app/models"

	"gorm.io/gorm"
)

func CreateFilmGenre(db *gorm.DB, filmID int64, genreID int64) *models.FilmGenre {
	filmGenre := &models.FilmGenre{
		FilmID:  filmID,
		GenreID: genreID,
	}

	db.Create(filmGenre)
	return filmGenre
}
