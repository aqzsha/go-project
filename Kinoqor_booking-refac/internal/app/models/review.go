package models

import "time"

func (Review) TableName() string {
	return "review"
}

type Review struct {
	ID              int64     `gorm:"primaryKey;column:id"`
	FilmID          int64
	UserID          int64 
	Body            string   
	Rating          float64   
	CreatedAt       *time.Time   
}
	