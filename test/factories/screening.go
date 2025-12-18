package factories

import (
	"time"

	"booking/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateScreening(
	db *gorm.DB,
	cinemaID int64,
	hallID int64,
	filmID int64,
) *models.Screening {
	startAt := time.Now().Add(time.Hour)
	endAt := startAt.Add(2 * time.Hour)

	priceInts, _ := faker.RandomInt(1000, 5000)
	price := 2000.0
	if len(priceInts) > 0 {
		price = float64(priceInts[0])
	}

	screening := &models.Screening{
		Start_At: startAt,
		End_At:   endAt,
		Language: faker.Word(),
		Format:   faker.Word(),
		Price:    price,
		CinemaID: cinemaID,
		HallID:   hallID,
		FilmID:   filmID,
	}

	db.Create(screening)
	return screening
}
