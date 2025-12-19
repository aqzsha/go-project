package factories

import (
	"time"

	"movies/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateCinemaDetails(db *gorm.DB) *models.CinemaDetails {
	name := faker.Word()
	description := faker.Sentence()
	address := faker.Word() + ", " + faker.Word()
	latitude := faker.Latitude()
	longitude := faker.Longitude()

	cinemaDetails := &models.CinemaDetails{
		Name:        name,
		Description: description,
		Address:     address,
		Latitude:    latitude,
		Longitude:   longitude,
		CreatedAt:   time.Now(),
	}

	db.Create(cinemaDetails)
	return cinemaDetails
}
