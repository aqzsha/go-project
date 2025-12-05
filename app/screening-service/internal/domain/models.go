package domain

import (
	"time"
)

type Screening struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	CinemaID   int       `json:"cinema_id" gorm:"not null"`
	HallID     int       `json:"hall_id" gorm:"not null"`
	FilmID     int       `json:"film_id" gorm:"not null"`
	StartAt    time.Time `json:"start_at" gorm:"not null"`
	EndAt      time.Time `json:"end_at" gorm:"not null"`
	Language   string    `json:"language" gorm:"size:50"`
	Format     string    `json:"format" gorm:"size:50"` // 2D, 3D, IMAX
	Price      float64   `json:"price" gorm:"type:decimal(10,2)"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Seat struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	HallID      int       `json:"hall_id" gorm:"not null"`
	ScreeningID int       `json:"screening_id" gorm:"not null"`
	Row         string    `json:"row" gorm:"size:10;not null"`
	Number      int       `json:"number" gorm:"not null"`
	Status      string    `json:"status" gorm:"size:20;default:'available'"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

const (
	SeatStatusAvailable = "available"
	SeatStatusReserved  = "reserved"
	SeatStatusSold      = "sold"
)

type CreateScreeningRequest struct {
	CinemaID int       `json:"cinema_id" binding:"required"`
	HallID   int       `json:"hall_id" binding:"required"`
	FilmID   int       `json:"film_id" binding:"required"`
	StartAt  time.Time `json:"start_at" binding:"required"`
	EndAt    time.Time `json:"end_at" binding:"required"`
	Language string    `json:"language" binding:"required"`
	Format   string    `json:"format" binding:"required"`
	Price    float64   `json:"price" binding:"required,gt=0"`
}

type CreateSeatRequest struct {
	HallID      int    `json:"hall_id" binding:"required"`
	ScreeningID int    `json:"screening_id" binding:"required"`
	Row         string `json:"row" binding:"required"`
	Number      int    `json:"number" binding:"required,gt=0"`
}

type ListScreeningsQuery struct {
	CinemaID  int       `form:"cinema_id"`
	FilmID    int       `form:"film_id"`
	HallID    int       `form:"hall_id"`
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
	Page      int       `form:"page" binding:"min=1"`
	PageSize  int       `form:"page_size" binding:"min=1,max=100"`
}

type ScreeningWithSeats struct {
	Screening      Screening            `json:"screening"`
	TotalSeats     int                  `json:"total_seats"`
	AvailableSeats int                  `json:"available_seats"`
	SoldSeats      int                  `json:"sold_seats"`
	ReservedSeats  int                  `json:"reserved_seats"`
	Seats          []Seat               `json:"seats,omitempty"`
	SeatMap        map[string][]SeatDTO `json:"seat_map,omitempty"` // grouped by row
}

type SeatDTO struct {
	ID     int    `json:"id"`
	Row    string `json:"row"`
	Number int    `json:"number"`
	Status string `json:"status"`
}

func (Screening) TableName() string {
	return "screening"
}

func (Seat) TableName() string {
	return "seat"
}