package factories

import (
	"time"

	"booking/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateFilmDetails(db *gorm.DB) *models.FilmDetails {
	RandNumber, _ := faker.RandomInt(1, 10)

	duration := faker.Word()
	production := faker.Word()
	director := faker.Name()
	rate := RandNumber[0]
	ageLimit := RandNumber[1]

	filmDetails := &models.FilmDetails{
		Duration:   duration,
		Premier:    time.Now().AddDate(-1, 0, 0),
		Production: production,
		Director:   director,
		Rate:       float64(rate),
		AgeLimit:   ageLimit,
	}

	db.Create(filmDetails)
	return filmDetails
}
