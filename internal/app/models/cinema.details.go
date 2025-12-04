package models

import "time"

func (CinemaDetails) TableName() string {
	return "cinema_details"
}

type CinemaDetails struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Address     string    `gorm:"column:address;type:varchar(255);not null" json:"address"`
	Latitude    float64   `gorm:"column:latitude;type:numeric(9,6);not null" json:"latitude"`
	Longitude   float64   `gorm:"column:longitude;type:numeric(9,6);not null" json:"longitude"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

