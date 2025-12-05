package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"screening-service/internal/domain"
	"screening-service/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrScreeningNotFound    = errors.New("screening not found")
	ErrHallNotAvailable     = errors.New("hall is not available for the specified time")
	ErrInvalidTimeRange     = errors.New("end time must be after start time")
	ErrScreeningInPast      = errors.New("screening cannot be scheduled in the past")
	ErrScreeningHasStarted  = errors.New("cannot modify screening that has already started")
)

type ScreeningService interface {
	CreateScreening(ctx context.Context, req domain.CreateScreeningRequest) (*domain.Screening, error)
	GetScreeningByID(ctx context.Context, id int) (*domain.ScreeningWithSeats, error)
	ListScreenings(ctx context.Context, query domain.ListScreeningsQuery) ([]domain.Screening, int64, error)
	DeleteScreening(ctx context.Context, id int) error
	GenerateSeatsForScreening(ctx context.Context, screeningID int, rows []string, seatsPerRow int) error
}

type screeningService struct {
	screeningRepo repository.ScreeningRepository
	seatRepo      repository.SeatRepository
}

func NewScreeningService(
	screeningRepo repository.ScreeningRepository,
	seatRepo repository.SeatRepository,
) ScreeningService {
	return &screeningService{
		screeningRepo: screeningRepo,
		seatRepo:      seatRepo,
	}
}

func (s *screeningService) CreateScreening(ctx context.Context, req domain.CreateScreeningRequest) (*domain.Screening, error) {
	if req.EndAt.Before(req.StartAt) || req.EndAt.Equal(req.StartAt) {
		return nil, ErrInvalidTimeRange
	}

	if req.StartAt.Before(time.Now()) {
		return nil, ErrScreeningInPast
	}

	available, err := s.screeningRepo.CheckHallAvailability(ctx, req.HallID, req.StartAt, req.EndAt, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check hall availability: %w", err)
	}
	if !available {
		return nil, ErrHallNotAvailable
	}

	screening := &domain.Screening{
		CinemaID: req.CinemaID,
		HallID:   req.HallID,
		FilmID:   req.FilmID,
		StartAt:  req.StartAt,
		EndAt:    req.EndAt,
		Language: req.Language,
		Format:   req.Format,
		Price:    req.Price,
	}

	if err := s.screeningRepo.Create(ctx, screening); err != nil {
		return nil, fmt.Errorf("failed to create screening: %w", err)
	}

	return screening, nil
}

func (s *screeningService) GetScreeningByID(ctx context.Context, id int) (*domain.ScreeningWithSeats, error) {
	screening, err := s.screeningRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScreeningNotFound
		}
		return nil, fmt.Errorf("failed to get screening: %w", err)
	}

	stats, err := s.seatRepo.GetSeatStatistics(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat statistics: %w", err)
	}

	seats, err := s.seatRepo.GetByScreeningID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get seats: %w", err)
	}

	seatMap := make(map[string][]domain.SeatDTO)
	for _, seat := range seats {
		seatMap[seat.Row] = append(seatMap[seat.Row], domain.SeatDTO{
			ID:     seat.ID,
			Row:    seat.Row,
			Number: seat.Number,
			Status: seat.Status,
		})
	}

	totalSeats := int(stats[domain.SeatStatusAvailable] + stats[domain.SeatStatusReserved] + stats[domain.SeatStatusSold])

	return &domain.ScreeningWithSeats{
		Screening:      *screening,
		TotalSeats:     totalSeats,
		AvailableSeats: int(stats[domain.SeatStatusAvailable]),
		SoldSeats:      int(stats[domain.SeatStatusSold]),
		ReservedSeats:  int(stats[domain.SeatStatusReserved]),
		SeatMap:        seatMap,
	}, nil
}

func (s *screeningService) ListScreenings(ctx context.Context, query domain.ListScreeningsQuery) ([]domain.Screening, int64, error) {
	screenings, total, err := s.screeningRepo.List(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list screenings: %w", err)
	}
	return screenings, total, nil
}

func (s *screeningService) DeleteScreening(ctx context.Context, id int) error {
	screening, err := s.screeningRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrScreeningNotFound
		}
		return fmt.Errorf("failed to get screening: %w", err)
	}

	if screening.StartAt.Before(time.Now()) {
		return ErrScreeningHasStarted
	}

	if err := s.screeningRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete screening: %w", err)
	}

	return nil
}

func (s *screeningService) GenerateSeatsForScreening(ctx context.Context, screeningID int, rows []string, seatsPerRow int) error {
	_, err := s.screeningRepo.GetByID(ctx, screeningID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrScreeningNotFound
		}
		return fmt.Errorf("failed to get screening: %w", err)
	}

	var seats []domain.Seat
	for _, row := range rows {
		for seatNum := 1; seatNum <= seatsPerRow; seatNum++ {
			seats = append(seats, domain.Seat{
				HallID:      0, // Should get from screening
				ScreeningID: screeningID,
				Row:         row,
				Number:      seatNum,
				Status:      domain.SeatStatusAvailable,
			})
		}
	}

	if err := s.seatRepo.CreateBatch(ctx, seats); err != nil {
		return fmt.Errorf("failed to create seats: %w", err)
	}

	return nil
}