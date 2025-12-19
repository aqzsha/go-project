package models

import "time"

type Ticket struct {
	ID              int64 `gorm:"primaryKey"`
	SeatScreeningID int64
	QrCode          string `gorm:"column:qr_code;uniqueIndex"`
	Status          string
	ScanAt          *time.Time `gorm:"column:scan_at"`
}

func (Ticket) TableName() string {
	return "ticket"
}
