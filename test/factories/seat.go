package factories

import (
	"booking/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateSeat(db *gorm.DB, hallID int64) *models.Seat {
	row, _ := faker.RandomInt(1, 20)
	number, _ := faker.RandomInt(1, 30)

	seat := &models.Seat{
		HallID: hallID,
		Row:    int64(row[0]),
		Number: int64(number[0]),
	}

	db.Create(seat)
	return seat
}
