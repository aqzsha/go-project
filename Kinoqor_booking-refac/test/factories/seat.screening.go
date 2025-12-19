package factories

import (
	"booking/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateSeatScreening(
	db *gorm.DB,
	hallID int64,
	screeningID int64,
) *models.SeatScreening {
	row, _ := faker.RandomInt(1, 20)
	number, _ := faker.RandomInt(1, 30)

	statuses := []string{"free", "reserved", "sold"}

	seatScreening := &models.SeatScreening{
		HallID:      hallID,
		ScreeningID: screeningID,
		Row:         int64(row[0]),
		Number:      int64(number[0]),
		Status:      statuses[0],
	}

	db.Create(seatScreening)
	return seatScreening
}
