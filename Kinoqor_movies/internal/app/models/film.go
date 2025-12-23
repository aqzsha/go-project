package models

import "time"

func (Film) TableName() string {
	return "film"
}

type Film struct {
	ID          int64       `json:"id" gorm:"primaryKey"`
	Name        string      `json:"name"`
	Description string      `json:"description"`

	DetailsID int64       `json:"details_id"`
	Details   FilmDetails `json:"details" gorm:"foreignKey:DetailsID;references:ID"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
}
