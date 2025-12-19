package factories

import (
	"booking/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateTicket(
	db *gorm.DB,
	seatScreeningID int64,
) *models.Ticket {
	qrCode := faker.UUIDDigit()

	ticket := &models.Ticket{
		SeatScreeningID: seatScreeningID,
		QrCode:          qrCode,
		Status:          "active",
		ScanAt:          nil,
	}

	db.Create(ticket)
	return ticket
}
