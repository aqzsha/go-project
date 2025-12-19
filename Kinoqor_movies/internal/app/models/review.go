package models

import "time"

func (Review) TableName() string {
	return "review"
}

type Review struct {
	ID        int64      `gorm:"primaryKey;column:id" json:"id"`
	FilmID    int64      `gorm:"column:film_id" json:"film_id"`
	UserID    int64      `gorm:"column:user_id" json:"user_id"`
	Title     string 	 `gorm:"column:title;type:text" json:"title"`
	Body      string     `gorm:"column:body;type:text" json:"body"`
	Rating    float64    `gorm:"column:rating" json:"rating"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
}
