package models

type Seat struct {
	ID     int64
	HallID int64
	Row    int64
	Number int64
}

func (Seat) TableName() string {
	return "seat"
}
