package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/repository"
	"booking-service/pkg/cache"
	"booking-service/pkg/qrcode"
	"booking-service/pkg/pdf"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTicketNotFound       = errors.New("ticket not found")
	ErrSeatsNotAvailable    = errors.New("one or more seats are not available")
	ErrReservationExpired   = errors.New("reservation has expired")
	ErrScreeningStarted     = errors.New("screening has already started")
	ErrScreeningNotFound    = errors.New("screening not found")
	ErrInvalidTicketStatus  = errors.New("invalid ticket status for this operation")
	ErrAlreadyCheckedIn     = errors.New("ticket already checked in")
	ErrCheckInNotAllowed    = errors.New("check-in not allowed at this time")
)

const (
	ReservationTTL = 5 * time.Minute
	CheckInWindow  = 30 * time.Minute //check-in 30 min before screening
)

type BookingService interface {
	ReserveSeats(ctx context.Context, req domain.ReserveSeatsRequest) (*domain.BookingResponse, error)
	ConfirmBooking(ctx context.Context, req domain.ConfirmBookingRequest) error
	CreateTicket(ctx context.Context, req domain.CreateTicketRequest) (*domain.BookingResponse, error)
	GetTicketByID(ctx context.Context, ticketID int) (*domain.Ticket, error)
	GetTicketByQRCode(ctx context.Context, qrCode string) (*domain.Ticket, error)
	GetUserTickets(ctx context.Context, userID int, query domain.ListTicketsQuery) ([]domain.Ticket, int64, error)
	CheckInTicket(ctx context.Context, qrCode string) (*domain.Ticket, error)
	CancelTicket(ctx context.Context, ticketID int, userID int) error
	ExpireReservations(ctx context.Context) error
	ExpireScreeningTickets(ctx context.Context) error
}

type bookingService struct {
	ticketRepo       repository.TicketRepository
	screeningClient  ScreeningClient
	cache            cache.RedisCache
	qrCodeGenerator  qrcode.Generator
	pdfGenerator     PDFGenerator
	notificationSvc  NotificationService
}

func NewBookingService(
	ticketRepo repository.TicketRepository,
	screeningClient ScreeningClient,
	cache cache.RedisCache,
	qrCodeGenerator qrcode.Generator,
	pdfGenerator PDFGenerator,
	notificationSvc NotificationService,
) BookingService {
	return &bookingService{
		ticketRepo:      ticketRepo,
		screeningClient: screeningClient,
		cache:           cache,
		qrCodeGenerator: qrCodeGenerator,
		pdfGenerator:    pdfGenerator,
		notificationSvc: notificationSvc,
	}
}

func (s *bookingService) ReserveSeats(ctx context.Context, req domain.ReserveSeatsRequest) (*domain.BookingResponse, error) {
	screening, err := s.screeningClient.GetScreening(ctx, req.ScreeningID)
	if err != nil {
		return nil, fmt.Errorf("failed to get screening: %w", err)
	}

	if screening.StartAt.Before(time.Now()) {
		return nil, ErrScreeningStarted
	}

	reservationID := uuid.New().String()
	expiresAt := time.Now().Add(ReservationTTL)

	for _, seatID := range req.SeatIDs {
		key := fmt.Sprintf("seat:reservation:%d:%d", req.ScreeningID, seatID)
		
		reservation := domain.SeatReservation{
			SeatID:      seatID,
			UserID:      req.UserID,
			ScreeningID: req.ScreeningID,
			ExpiresAt:   expiresAt,
		}
		
		data, _ := json.Marshal(reservation)
		
		success, err := s.cache.SetNX(ctx, key, string(data), ReservationTTL)
		if err != nil || !success {
			s.releaseSeats(ctx, req.ScreeningID, req.SeatIDs)
			return nil, ErrSeatsNotAvailable
		}
		
		if err := s.screeningClient.UpdateSeatStatus(ctx, seatID, "reserved"); err != nil {
			s.releaseSeats(ctx, req.ScreeningID, req.SeatIDs)
			return nil, fmt.Errorf("failed to update seat status: %w", err)
		}
	}

	var tickets []domain.Ticket
	totalPrice := 0.0

	for _, seatID := range req.SeatIDs {
		qrCode := s.generateUniqueQRCode()
		
		ticket := domain.Ticket{
			UserID:      req.UserID,
			ScreeningID: req.ScreeningID,
			SeatID:      seatID,
			Price:       screening.Price,
			Status:      domain.TicketStatusReserved,
			QRCode:      qrCode,
			ExpiresAt:   expiresAt,
		}
		tickets = append(tickets, ticket)
		totalPrice += screening.Price
	}

	if err := s.ticketRepo.CreateBatch(ctx, tickets); err != nil {
		s.releaseSeats(ctx, req.ScreeningID, req.SeatIDs)
		return nil, fmt.Errorf("failed to create tickets: %w", err)
	}

	return &domain.BookingResponse{
		Tickets:       tickets,
		TotalPrice:    totalPrice,
		ExpiresAt:     expiresAt,
		ReservationID: reservationID,
	}, nil
}

func (s *bookingService) ConfirmBooking(ctx context.Context, req domain.ConfirmBookingRequest) error {
	tickets, err := s.ticketRepo.GetByIDs(ctx, req.TicketIDs)
	if err != nil {
		return fmt.Errorf("failed to get tickets: %w", err)
	}

	if len(tickets) == 0 {
		return ErrTicketNotFound
	}

	for _, ticket := range tickets {
		if ticket.UserID != req.UserID {
			return errors.New("unauthorized: ticket does not belong to user")
		}
		if ticket.Status != domain.TicketStatusReserved {
			return ErrInvalidTicketStatus
		}
		if time.Now().After(ticket.ExpiresAt) {
			return ErrReservationExpired
		}
	}

	if err := s.ticketRepo.UpdateStatusBatch(ctx, req.TicketIDs, domain.TicketStatusPaid); err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}

	var seatIDs []int
	for _, ticket := range tickets {
		seatIDs = append(seatIDs, ticket.SeatID)
	}

	for _, seatID := range seatIDs {
		if err := s.screeningClient.UpdateSeatStatus(ctx, seatID, "sold"); err != nil {
			return fmt.Errorf("failed to update seat status: %w", err)
		}
	}

	s.releaseSeats(ctx, tickets[0].ScreeningID, seatIDs)

	//generator QR codes, PDFs and send emails for each ticket
	for i := range tickets {
		tickets[i].Status = domain.TicketStatusPaid
		
		ticketDetails, err := s.getTicketDetails(ctx, &tickets[i])
		if err != nil {
			log.Printf("Error getting ticket details: %v", err)
			continue
		}
		
		qrCodeBase64, err := s.qrCodeGenerator.Generate(tickets[i].QRCode)
		if err != nil {
			log.Printf("Error generating QR code: %v", err)
			continue
		}
		
		pdfData, err := s.pdfGenerator.GenerateTicket(ctx, &tickets[i], ticketDetails, qrCodeBase64)
		if err != nil {
			log.Printf("Error generating PDF: %v", err)
			continue
		}
		
		//upload PDF to MinIO
		pdfPath, err := s.storageService.UploadTicket(ctx, tickets[i].ID, pdfData)
		if err != nil {
			log.Printf("Error uploading PDF: %v", err)
		} else {
			tickets[i].PDFPath = pdfPath
			s.ticketRepo.Update(ctx, &tickets[i])
		}
		
		//generate download URL
		downloadURL := ""
		if pdfPath != "" {
			downloadURL, _ = s.storageService.GetTicketURL(ctx, pdfPath, 7*24*time.Hour)
		}
		
		//get user email
		userEmail, err := s.getUserEmail(ctx, req.UserID)
		if err == nil {
			//send email with PDF attachment
			go s.notificationSvc.SendTicketEmail(context.Background(), userEmail, ticketDetails, pdfData, downloadURL)
		}
	}

	//send overall confirmation notification
	s.notificationSvc.SendBookingConfirmation(ctx, req.UserID, tickets)

	return nil
}

func (s *bookingService) CreateTicket(ctx context.Context, req domain.CreateTicketRequest) (*domain.BookingResponse, error) {
	//direct purchase flow (reserve + confirm in one step)
	reserveReq := domain.ReserveSeatsRequest{
		UserID:      req.UserID,
		ScreeningID: req.ScreeningID,
		SeatIDs:     req.SeatIDs,
	}

	booking, err := s.ReserveSeats(ctx, reserveReq)
	if err != nil {
		return nil, err
	}

	//auto-confirm (assuming payment is handled externally)
	var ticketIDs []int
	for _, ticket := range booking.Tickets {
		ticketIDs = append(ticketIDs, ticket.ID)
	}

	confirmReq := domain.ConfirmBookingRequest{
		UserID:      req.UserID,
		TicketIDs:   ticketIDs,
		PaymentID:   uuid.New().String(),
		TotalAmount: booking.TotalPrice,
	}

	if err := s.ConfirmBooking(ctx, confirmReq); err != nil {
		return nil, err
	}

	updatedTickets, _ := s.ticketRepo.GetByIDs(ctx, ticketIDs)
	booking.Tickets = updatedTickets

	return booking, nil
}

func (s *bookingService) GetTicketByID(ctx context.Context, ticketID int) (*domain.Ticket, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	return ticket, nil
}

func (s *bookingService) GetTicketByQRCode(ctx context.Context, qrCode string) (*domain.Ticket, error) {
	ticket, err := s.ticketRepo.GetByQRCode(ctx, qrCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	return ticket, nil
}

func (s *bookingService) GetUserTickets(ctx context.Context, userID int, query domain.ListTicketsQuery) ([]domain.Ticket, int64, error) {
	return s.ticketRepo.GetByUserID(ctx, userID, query)
}

func (s *bookingService) CheckInTicket(ctx context.Context, qrCode string) (*domain.Ticket, error) {
	ticket, err := s.GetTicketByQRCode(ctx, qrCode)
	if err != nil {
		return nil, err
	}

	if ticket.Status == domain.TicketStatusCheckedIn {
		return nil, ErrAlreadyCheckedIn
	}

	if ticket.Status != domain.TicketStatusPaid {
		return nil, ErrInvalidTicketStatus
	}

	screening, err := s.screeningClient.GetScreening(ctx, ticket.ScreeningID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	allowCheckInFrom := screening.StartAt.Add(-CheckInWindow)
	
	if now.Before(allowCheckInFrom) {
		return nil, fmt.Errorf("check-in opens %d minutes before screening", int(CheckInWindow.Minutes()))
	}

	if now.After(screening.EndAt) {
		return nil, errors.New("screening has ended")
	}

	if err := s.ticketRepo.CheckIn(ctx, ticket.ID, now); err != nil {
		return nil, err
	}

	ticket.Status = domain.TicketStatusCheckedIn
	ticket.CheckedInAt = &now

	return ticket, nil
}

func (s *bookingService) CancelTicket(ctx context.Context, ticketID int, userID int) error {
	ticket, err := s.GetTicketByID(ctx, ticketID)
	if err != nil {
		return err
	}

	if ticket.UserID != userID {
		return errors.New("unauthorized")
	}

	if ticket.Status != domain.TicketStatusPaid && ticket.Status != domain.TicketStatusReserved {
		return ErrInvalidTicketStatus
	}

	if err := s.ticketRepo.UpdateStatus(ctx, ticketID, domain.TicketStatusCancelled); err != nil {
		return err
	}

	s.screeningClient.UpdateSeatStatus(ctx, ticket.SeatID, "available")

	return nil
}

func (s *bookingService) ExpireReservations(ctx context.Context) error {
	tickets, err := s.ticketRepo.GetExpiredReservations(ctx, time.Now())
	if err != nil {
		return err
	}

	for _, ticket := range tickets {
		s.ticketRepo.UpdateStatus(ctx, ticket.ID, domain.TicketStatusExpired)
		s.screeningClient.UpdateSeatStatus(ctx, ticket.SeatID, "available")
		
		s.notificationSvc.SendReservationExpired(ctx, ticket.UserID, ticket.ID)
	}

	return nil
}

func (s *bookingService) ExpireScreeningTickets(ctx context.Context) error {
	tickets, err := s.ticketRepo.GetTicketsForExpiredScreenings(ctx, time.Now())
	if err != nil {
		return err
	}

	for _, ticket := range tickets {
		if ticket.Status != domain.TicketStatusCheckedIn {
			s.ticketRepo.UpdateStatus(ctx, ticket.ID, domain.TicketStatusExpired)
			
			//if user missed the screening
			s.notificationSvc.SendScreeningMissed(ctx, ticket.UserID, ticket.ID)
		}
	}

	return nil
}

func (s *bookingService) releaseSeats(ctx context.Context, screeningID int, seatIDs []int) {
	for _, seatID := range seatIDs {
		key := fmt.Sprintf("seat:reservation:%d:%d", screeningID, seatID)
		s.cache.Del(ctx, key)
		s.screeningClient.UpdateSeatStatus(ctx, seatID, "available")
	}
}

func (s *bookingService) generateUniqueQRCode() string {
	return fmt.Sprintf("TKT-%s", uuid.New().String())
}

type ScreeningClient interface {
	GetScreening(ctx context.Context, screeningID int) (*domain.ScreeningInfo, error)
	UpdateSeatStatus(ctx context.Context, seatID int, status string) error
}

type PDFGenerator interface {
	GenerateTicket(ctx context.Context, ticket *domain.Ticket) (string, error)
}

type NotificationService interface {
	SendBookingConfirmation(ctx context.Context, userID int, tickets []domain.Ticket) error
	SendReservationExpired(ctx context.Context, userID int, ticketID int) error
	SendScreeningMissed(ctx context.Context, userID int, ticketID int) error
}