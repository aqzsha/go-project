package domain

import (
	"time"
)

//ticket
type Ticket struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	UserID      int       `json:"user_id" gorm:"not null;index"`
	ScreeningID int       `json:"screening_id" gorm:"not null;index"`
	SeatID      int       `json:"seat_id" gorm:"not null;uniqueIndex"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2)"`
	Status      string    `json:"status" gorm:"size:20;not null;index"` // reserved, paid, checked_in, expired, cancelled
	QRCode      string    `json:"qr_code" gorm:"size:255"`
	PDFPath     string    `json:"pdf_path" gorm:"size:500"`
	ExpiresAt   time.Time `json:"expires_at"`
	CheckedInAt *time.Time `json:"checked_in_at"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

const (
	TicketStatusReserved   = "reserved"   // Seat temporarily held (5 min)
	TicketStatusPaid       = "paid"       // Payment completed
	TicketStatusCheckedIn  = "checked_in" // User checked in at cinema
	TicketStatusExpired    = "expired"    // Reservation expired or screening passed
	TicketStatusCancelled  = "cancelled"  // Cancelled by user or system
)

type SeatReservation struct {
	SeatID      int       `json:"seat_id"`
	UserID      int       `json:"user_id"`
	ScreeningID int       `json:"screening_id"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type CreateTicketRequest struct {
	UserID      int   `json:"user_id" binding:"required"`
	ScreeningID int   `json:"screening_id" binding:"required"`
	SeatIDs     []int `json:"seat_ids" binding:"required,min=1,max=10"`
}

type ReserveSeatsRequest struct {
	UserID      int   `json:"user_id" binding:"required"`
	ScreeningID int   `json:"screening_id" binding:"required"`
	SeatIDs     []int `json:"seat_ids" binding:"required,min=1,max=10"`
}

type ConfirmBookingRequest struct {
	UserID       int     `json:"user_id" binding:"required"`
	TicketIDs    []int   `json:"ticket_ids" binding:"required"`
	PaymentID    string  `json:"payment_id" binding:"required"`
	TotalAmount  float64 `json:"total_amount" binding:"required,gt=0"`
}

type CheckInRequest struct {
	QRCode string `json:"qr_code" binding:"required"`
}

type TicketWithDetails struct {
	Ticket      Ticket           `json:"ticket"`
	Screening   ScreeningInfo    `json:"screening"`
	Seat        SeatInfo         `json:"seat"`
	Cinema      CinemaInfo       `json:"cinema"`
	Film        FilmInfo         `json:"film"`
}

type ScreeningInfo struct {
	ID       int       `json:"id"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
	Language string    `json:"language"`
	Format   string    `json:"format"`
}

type SeatInfo struct {
	ID     int    `json:"id"`
	Row    string `json:"row"`
	Number int    `json:"number"`
}

type CinemaInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type FilmInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Poster string `json:"poster,omitempty"`
}

type BookingResponse struct {
	Tickets      []Ticket  `json:"tickets"`
	TotalPrice   float64   `json:"total_price"`
	ExpiresAt    time.Time `json:"expires_at"`
	ReservationID string   `json:"reservation_id"`
}

type ListTicketsQuery struct {
	UserID      int       `form:"user_id"`
	ScreeningID int       `form:"screening_id"`
	Status      string    `form:"status"`
	StartDate   time.Time `form:"start_date"`
	EndDate     time.Time `form:"end_date"`
	Page        int       `form:"page" binding:"min=1"`
	PageSize    int       `form:"page_size" binding:"min=1,max=100"`
}

func (Ticket) TableName() string {
	return "ticket"
}