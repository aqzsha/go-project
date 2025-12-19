package models

func (Hall) TableName() string {
	return "hall"
}

type Hall struct {
	ID       int64  `gorm:"primaryKey;column:id" json:"id"`
	CinemaID int64  `gorm:"column:cinema_id" json:"cinema_id"`
	Name     string `gorm:"column:name" json:"name"`
	Seats    int64  `gorm:"column:seats" json:"seats"`
}