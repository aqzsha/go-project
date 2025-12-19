package models

import "time"

func (Film) TableName() string {
	return "film"
}

type Film struct {
	ID          int64       `json:"id" gorm:"primaryKey;column:id"`
	Name        string      `json:"name" gorm:"column:name;not null"`
	Description string      `json:"description" gorm:"column:description"`
	DetailsID   int64       `json:"details_id" gorm:"column:details_id;not null"`
	Details     FilmDetails `json:"details"`
	StartDate   time.Time   `json:"start_date"`
	EndDate     time.Time   `json:"end_date"`
	CreatedAt   time.Time   `json:"created_at"`
}
