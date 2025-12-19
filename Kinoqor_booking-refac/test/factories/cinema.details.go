package factories

import (
	"time"

	"booking/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateCinemaDetails(db *gorm.DB) *models.CinemaDetails {
	cinemaDetails := &models.CinemaDetails{
		Name:        faker.Word(),
		Description: faker.Sentence(),
		Address:     faker.Word() + ", " + faker.Word(),
		Latitude:    faker.Latitude(),
		Longitude:   faker.Longitude(),
		CreatedAt:   time.Now(),
	}

	db.Create(cinemaDetails)
	return cinemaDetails
}
