package models

import "time"

type Screening struct {
	ID       int64
	Start_At time.Time
	End_At   time.Time
	Language string
	Format   string
	Price    float64
	HallID   int64
	FilmID   int64
	CinemaID int64
}

func (Screening) TableName() string {
	return "screening"
}
