package factories

import (
	"movies/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateCinema(db *gorm.DB, details *models.CinemaDetails) *models.Cinema {
	name := faker.Word()

	cinema := &models.Cinema{
		Name:      name,
		DetailsID: details.ID, 
	}

	db.Create(cinema)
	return cinema
}
