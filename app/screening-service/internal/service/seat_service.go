package service

import (
	"context"
	"errors"
	"fmt"

	"screening-service/internal/domain"
	"screening-service/internal/repository"
	"screening-service/pkg/cache"

	"gorm.io/gorm"
)

var (
	ErrSeatNotFound        = errors.New("seat not found")
	ErrSeatAlreadyExists   = errors.New("seat already exists for this screening")
	ErrSeatOccupied        = errors.New("seat is occupied and cannot be deleted")
	ErrInvalidSeatStatus   = errors.New("invalid seat status")
)

type SeatService interface {
	CreateSeat(ctx context.Context, req domain.CreateSeatRequest) (*domain.Seat, error)
	GetSeatByID(ctx context.Context, id int) (*domain.Seat, error)
	GetScreeningSeats(ctx context.Context, screeningID int) ([]domain.Seat, error)
	UpdateSeatStatus(ctx context.Context, seatID int, status string) error
	DeleteSeat(ctx context.Context, id int) error
	CheckSeatsAvailability(ctx context.Context, seatIDs []int) (bool, error)
}

type seatService struct {
	seatRepo      repository.SeatRepository
	screeningRepo repository.ScreeningRepository
	cache         cache.RedisCache
}

func NewSeatService(
	seatRepo repository.SeatRepository,
	screeningRepo repository.ScreeningRepository,
	cache cache.RedisCache,
) SeatService {
	return &seatService{
		seatRepo:      seatRepo,
		screeningRepo: screeningRepo,
		cache:         cache,
	}
}

func (s *seatService) CreateSeat(ctx context.Context, req domain.CreateSeatRequest) (*domain.Seat, error) {
	_, err := s.screeningRepo.GetByID(ctx, req.ScreeningID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScreeningNotFound
		}
		return nil, fmt.Errorf("failed to verify screening: %w", err)
	}

	seat := &domain.Seat{
		HallID:      req.HallID,
		ScreeningID: req.ScreeningID,
		Row:         req.Row,
		Number:      req.Number,
		Status:      domain.SeatStatusAvailable,
	}

	if err := s.seatRepo.Create(ctx, seat); err != nil {
		return nil, fmt.Errorf("failed to create seat: %w", err)
	}

	return seat, nil
}

func (s *seatService) GetSeatByID(ctx context.Context, id int) (*domain.Seat, error) {
	seat, err := s.seatRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSeatNotFound
		}
		return nil, fmt.Errorf("failed to get seat: %w", err)
	}
	return seat, nil
}

func (s *seatService) GetScreeningSeats(ctx context.Context, screeningID int) ([]domain.Seat, error) {
	seats, err := s.seatRepo.GetByScreeningID(ctx, screeningID)
	if err != nil {
		return nil, fmt.Errorf("failed to get screening seats: %w", err)
	}
	return seats, nil
}

func (s *seatService) UpdateSeatStatus(ctx context.Context, seatID int, status string) error {
	validStatuses := map[string]bool{
		domain.SeatStatusAvailable: true,
		domain.SeatStatusReserved:  true,
		domain.SeatStatusSold:      true,
	}

	if !validStatuses[status] {
		return ErrInvalidSeatStatus
	}

	_, err := s.GetSeatByID(ctx, seatID)
	if err != nil {
		return err
	}

	if err := s.seatRepo.UpdateStatus(ctx, seatID, status); err != nil {
		return fmt.Errorf("failed to update seat status: %w", err)
	}

	return nil
}

func (s *seatService) DeleteSeat(ctx context.Context, id int) error {
	seat, err := s.GetSeatByID(ctx, id)
	if err != nil {
		return err
	}

	if seat.Status != domain.SeatStatusAvailable {
		return ErrSeatOccupied
	}

	if err := s.seatRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete seat: %w", err)
	}

	return nil
}

func (s *seatService) CheckSeatsAvailability(ctx context.Context, seatIDs []int) (bool, error) {
	available, err := s.seatRepo.CheckSeatsAvailability(ctx, seatIDs)
	if err != nil {
		return false, fmt.Errorf("failed to check seat availability: %w", err)
	}
	return available, nil
}