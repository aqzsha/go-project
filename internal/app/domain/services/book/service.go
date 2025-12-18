package book

import (
	"booking/internal/app/core/helpers/qrCode"
	"booking/internal/app/core/usecase"
	repository "booking/internal/app/domain/repositories/book"
	"booking/internal/app/models"
	"booking/pkg/cache"
	"booking/pkg/gormtx"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound        = errors.New("seatScreening not found")
	ErrAlreadyReserved = errors.New("seatScreening is already reserved")
	ErrBookingExpired  = errors.New("seatScreening is already expired")
)

const (
	SeatStatusFree = "free"
	SeatStatusSold = "sold"

	TicketStatusActive = "active"
	TicketStatusUsed = "used"
)

type Service interface {
	Create(ctx context.Context, screeningSeatID int64) (bool, error)
	Get(ctx context.Context, ticketID int64) (models.Ticket, error)
	Delete(ctx context.Context, ticketID int64) (bool, error)
	PaidTicket(ctx context.Context, screeningSeatID int64) ([]byte, error)
	ScanTicket(ctx context.Context, token string) error
}

type service struct {
	db         *gorm.DB
	cache      cache.Cache
	tx         usecase.TxManager
	repository repository.Repository
}

func NewService(
	db *gorm.DB,
	cache cache.Cache,
	tx usecase.TxManager,
) Service {
	return &service{
		db:         db,
		repository: repository.NewRepository(db),
		cache:      cache,
		tx:         tx,
	}
}

const bookingTTL = 10 * time.Minute

func (s *service) Create(ctx context.Context, screeningSeatID int64) (bool, error) {
	key := fmt.Sprintf("seat:%d", screeningSeatID)

	ok, err := s.cache.SetNX(
		ctx,
		key,
		[]byte("reserved"),
		bookingTTL,
	)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, ErrAlreadyReserved
	}

	ok, err = s.repository.Create(ctx, screeningSeatID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *service) Get(ctx context.Context, ticketID int64) (models.Ticket, error) {
	ticket, err := s.repository.Get(ctx, ticketID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return models.Ticket{}, ErrNotFound
		}
	}

	return ticket, err
}

func (s *service) Delete(ctx context.Context, ticketID int64) (bool, error) {
	ok, err := s.repository.Delete(ctx, ticketID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ok, ErrNotFound
		}
	}

	return ok, err
}

func (s *service) PaidTicket(ctx context.Context, screeningSeatID int64) ([]byte, error) {
	key := fmt.Sprintf("seat:%d", screeningSeatID)

	if _, err := s.cache.Get(ctx, key); err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			return nil, ErrBookingExpired
		}
		return nil, err
	}

	var (
		ticket    models.Ticket
		seat      models.SeatScreening
		screening models.Screening
		film      models.Film
		hall      models.Hall
		cinema    models.Cinema
	)

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		tx := gormtx.From(ctx)
		repo := repository.NewRepository(tx)

		if err := repo.MarkSeatAsSold(ctx, screeningSeatID); err != nil {
			return err
		}

		token := uuid.NewString()

		ticket = models.Ticket{
			SeatScreeningID: screeningSeatID,
			QrCode:          token,
			Status:          TicketStatusActive,
		}
		if err := repo.CreateTicket(ctx, &ticket); err != nil {
			return err
		}

		if err := repo.GetSeat(ctx, &seat, screeningSeatID); err != nil {
			return err
		}
		if err := repo.GetScreening(ctx, &screening, seat.ScreeningID); err != nil {
			return err
		}
		if err := repo.GetFilm(ctx, &film, screening.FilmID); err != nil {
			return err
		}
		if err := repo.GetHall(ctx, &hall, screening.HallID); err != nil {
			return err
		}
		if err := repo.GetCinema(ctx, &cinema, screening.CinemaID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = s.cache.Delete(ctx, key)

	scanURL := fmt.Sprintf(
		"http://127.0.0.1:8881/tickets/scan/%s",
		ticket.QrCode,
	)

	qrPNG, err := qrCode.GenerateQRCode(scanURL)
	if err != nil {
		return nil, err
	}

	pdfBytes, err := qrCode.GenerateTicketPDF(
		ticket,
		screening,
		seat,
		film,
		hall,
		cinema,
		qrPNG,
	)
	if err != nil {
		return nil, err
	}

	return pdfBytes, nil
}

func (s *service) ScanTicket(ctx context.Context, token string) error {
	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		tx := gormtx.From(ctx)
		repo := repository.NewRepository(tx)

		ticket, err := repo.GetTicketByToken(ctx, token)
		if err != nil {
			return errors.New("билет не найден")
		}

		if ticket.Status == TicketStatusUsed {
			return errors.New("билет уже использован")
		}

		now := time.Now()

		if err := repo.MarkTicketUsed(ctx, ticket.ID, &now); err != nil {
			return err
		}

		if err := repo.MarkSeatUsed(ctx, ticket.SeatScreeningID); err != nil {
			return err
		}

		return nil
	})
}
