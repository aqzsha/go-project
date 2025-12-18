package models

import "time"

func (Film) TableName() string {
	return "film"
}

type Film struct {
	ID          int64        `gorm:"primaryKey;column:id"`
	Name        string       `gorm:"column:name;type:varchar(255);not null"`
	Description string       `gorm:"column:description;type:text"`
	DetailsID   int64        `gorm:"column:details_id;not null"`                     
	Details     FilmDetails  `gorm:"foreignKey:DetailsID;references:ID"`             
	StartDate   time.Time    `gorm:"column:start_date;not null"`
	EndDate     time.Time    `gorm:"column:end_date;not null"`
	CreatedAt   time.Time    `gorm:"column:created_at;autoCreateTime"`
}
