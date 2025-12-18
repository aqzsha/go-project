package book

import (
	"booking/internal/app/models"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("seatScreening not found")
var ErrAlreadyReserved = errors.New("seatScreening is already reserved")
var ErrAlreadySold = errors.New("seatScreening is already sold")

const (
	TicketStatusActive = "active"
	TicketStatusUsed   = "used"

	SeatStatusFree = "free"
	SeatStatusSold = "sold"
	SeatStatusUsed = "used"
)

type Repository interface {
	Create(ctx context.Context, seatScreeningID int64) (bool, error)
	Get(ctx context.Context, ticketID int64) (models.Ticket, error)
	Delete(ctx context.Context, ticketID int64) (bool, error)
	MarkAsSold(ctx context.Context, screeningSeatID int64) error
	CreateTicket(ctx context.Context, t *models.Ticket) error
	MarkSeatAsSold(ctx context.Context, seatID int64) error
	GetSeat(ctx context.Context, s *models.SeatScreening, id int64) error
	GetScreening(ctx context.Context, s *models.Screening, id int64) error
	GetFilm(ctx context.Context, f *models.Film, id int64) error
	GetHall(ctx context.Context, h *models.Hall, id int64) error
	GetCinema(ctx context.Context, c *models.Cinema, id int64) error
	MarkTicketUsed(
		ctx context.Context,
		ticketID int64,
		scanAt *time.Time,
	) error
	MarkSeatUsed(
		ctx context.Context,
		seatID int64,
	) error
	GetTicketByToken(
		ctx context.Context,
		token string,
	) (models.Ticket, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(
	ctx context.Context,
	seatScreeningID int64,
) (bool, error) {
	res := r.db.WithContext(ctx).
		Model(&models.SeatScreening{}).
		Where("id = ? AND status = ?", seatScreeningID, "free").
		Update("status", "reserved")

	if res.Error != nil {
		return false, res.Error
	}

	if res.RowsAffected == 0 {
		return false, ErrAlreadyReserved
	}

	return true, nil
}

func (r *repository) Get(ctx context.Context, ticketID int64) (models.Ticket, error) {
	var ticket models.Ticket

	if err := r.db.WithContext(ctx).
		Where("id = ?", ticketID).Where("status = ?", "available").
		First(&ticket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ticket, ErrNotFound
		}

		return ticket, fmt.Errorf("failed to get Ticket: %w", err)
	}

	return ticket, nil
}

func (r *repository) Delete(ctx context.Context, ticketID int64) (bool, error) {
	res := r.db.WithContext(ctx).Where("id = ?", ticketID).Delete(&models.Ticket{})
	if res.Error != nil {
		return false, fmt.Errorf("failed to delete Ticket: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return false, ErrNotFound
	}

	return true, nil
}

func (r *repository) MarkAsSold(
	ctx context.Context,
	screeningSeatID int64,
) error {
	res := r.db.WithContext(ctx).
		Model(&models.SeatScreening{}).
		Where("id = ? AND status != ?", screeningSeatID, "sold").
		Update("status", "sold")

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrAlreadySold
	}

	return nil
}

func (r *repository) MarkSeatAsSold(ctx context.Context, seatID int64) error {
	res := r.db.WithContext(ctx).
		Model(&models.SeatScreening{}).
		Where("id = ? AND status != ?", seatID, "sold").
		Update("status", "sold")

	if res.RowsAffected == 0 {
		return ErrAlreadySold
	}
	return res.Error
}

func (r *repository) CreateTicket(ctx context.Context, t *models.Ticket) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *repository) GetSeat(ctx context.Context, s *models.SeatScreening, id int64) error {
	return r.db.WithContext(ctx).First(s, id).Error
}

func (r *repository) GetScreening(ctx context.Context, s *models.Screening, id int64) error {
	return r.db.WithContext(ctx).First(s, id).Error
}

func (r *repository) GetFilm(ctx context.Context, f *models.Film, id int64) error {
	return r.db.WithContext(ctx).Preload("Details").First(f, id).Error
}

func (r *repository) GetHall(ctx context.Context, h *models.Hall, id int64) error {
	return r.db.WithContext(ctx).First(h, id).Error
}

func (r *repository) GetCinema(ctx context.Context, c *models.Cinema, id int64) error {
	return r.db.WithContext(ctx).Preload("Details").First(c, id).Error
}

func (r *repository) GetTicketByToken(
	ctx context.Context,
	token string,
) (models.Ticket, error) {

	var ticket models.Ticket

	err := r.db.WithContext(ctx).
		Where("qr_code = ?", token).
		First(&ticket).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Ticket{}, ErrNotFound
		}
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (r *repository) MarkTicketUsed(
	ctx context.Context,
	ticketID int64,
	scanAt *time.Time,
) error {

	return r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Where("id = ?", ticketID).
		Updates(map[string]any{
			"status":  TicketStatusUsed,
			"scan_at": scanAt,
		}).Error
}

func (r *repository) MarkSeatUsed(
	ctx context.Context,
	seatID int64,
) error {

	return r.db.WithContext(ctx).
		Model(&models.SeatScreening{}).
		Where("id = ?", seatID).
		Update("status", SeatStatusUsed).Error
}
