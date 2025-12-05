package worker

import (
	"context"
	"log"
	"time"

	"booking-service/internal/service"
)

type TicketWorker struct {
	bookingService service.BookingService
	stopChan       chan struct{}
}

func NewTicketWorker(bookingService service.BookingService) *TicketWorker {
	return &TicketWorker{
		bookingService: bookingService,
		stopChan:       make(chan struct{}),
	}
}

func (w *TicketWorker) Start(ctx context.Context) {
	reservationTicker := time.NewTicker(1 * time.Minute)
	defer reservationTicker.Stop()

	screeningTicker := time.NewTicker(15 * time.Minute)
	defer screeningTicker.Stop()

	log.Println("Ticket worker started")

	for {
		select {
		case <-reservationTicker.C:
			w.expireReservations(ctx)
		case <-screeningTicker.C:
			w.expireScreeningTickets(ctx)
		case <-w.stopChan:
			log.Println("Ticket worker stopped")
			return
		case <-ctx.Done():
			log.Println("Ticket worker context cancelled")
			return
		}
	}
}

func (w *TicketWorker) Stop() {
	close(w.stopChan)
}

func (w *TicketWorker) expireReservations(ctx context.Context) {
	log.Println("Running reservation expiration check...")
	
	if err := w.bookingService.ExpireReservations(ctx); err != nil {
		log.Printf("Error expiring reservations: %v", err)
		return
	}
	
	log.Println("Reservation expiration check completed")
}

func (w *TicketWorker) expireScreeningTickets(ctx context.Context) {
	log.Println("Running screening ticket expiration check...")
	
	if err := w.bookingService.ExpireScreeningTickets(ctx); err != nil {
		log.Printf("Error expiring screening tickets: %v", err)
		return
	}
	
	log.Println("Screening ticket expiration check completed")
}