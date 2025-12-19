package models

type SeatScreening struct {
	ID          int64
	HallID      int64
	ScreeningID int64
	Row         int64
	Number      int64
	Status      string
}

func (SeatScreening) TableName() string {
	return "seat_screening"
}
