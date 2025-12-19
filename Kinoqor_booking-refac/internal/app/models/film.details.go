package models

import "time"

func (FilmDetails) TableName() string {
	return "film_details"
}

type FilmDetails struct {
	ID         int64     `gorm:"primaryKey;column:id"`
	Duration   string    `gorm:"column:duration;type:text"`
	Premier    time.Time `gorm:"column:premier;not null"`
	Production string    `gorm:"column:production;type:varchar(255)"`
	Director   string    `gorm:"column:director;type:varchar(255)"`
	Rate       float64   `gorm:"column:rate"`
	AgeLimit   int       `gorm:"column:age_limit"`
}
	