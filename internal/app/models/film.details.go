package models

import "time"

func (FilmDetails) TableName() string {
	return "film_details"
}

type FilmDetails struct {
	ID         int64     `json:"id" gorm:"primaryKey;column:id"`
	Duration   string    `json:"duration" gorm:"column:duration;type:text"`
	Premier    time.Time `json:"premier" gorm:"column:premier;not null"`
	Production string    `json:"production" gorm:"column:production;type:varchar(255)"`
	Director   string    `json:"director" gorm:"column:director;type:varchar(255)"`
	Rate       float64   `json:"rate" gorm:"column:rate"`
	AgeLimit   int       `json:"age_limit" gorm:"column:age_limit"`
}
