package models

func (Hall) TableName() string {
	return "hall"
}

type Hall struct {
	ID       int64 `gorm:"primaryKey;column:id"`
	CinemaID int64
	Name     string
	Seats    int64
}