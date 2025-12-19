package factories

import (
	"movies/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateHall(db *gorm.DB, cinemaID int64) *models.Hall {
	name := faker.Word()
	seats, _ := faker.RandomInt(50, 300)

	hall := &models.Hall{
		CinemaID: cinemaID,
		Name:     name,
		Seats:    int64(seats[0]),
	}

	db.Create(hall)
	return hall
}
