package factories

import (
	"time"

	"booking/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateFilm(db *gorm.DB, details *models.FilmDetails) *models.Film {
	name := faker.Word()
	description := faker.Sentence()

	film := &models.Film{
		Name:        name,
		Description: description,
		DetailsID:   details.ID,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
	}

	db.Create(film)
	return film
}
