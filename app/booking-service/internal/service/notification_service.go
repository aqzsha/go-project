package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"booking-service/internal/domain"
	"booking-service/pkg/email"
	"booking-service/pkg/storage"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationService interface {
	SendBookingConfirmation(ctx context.Context, userID int, tickets []domain.Ticket) error
	SendReservationExpired(ctx context.Context, userID int, ticketID int) error
	SendScreeningMissed(ctx context.Context, userID int, ticketID int) error
	SendTicketEmail(ctx context.Context, userEmail string, ticket *domain.TicketWithDetails, pdfData []byte, downloadURL string) error
}

type notificationService struct {
	emailService    email.EmailService
	storageService  storage.MinIOStorage
	rabbitConn      *amqp.Connection
	rabbitChannel   *amqp.Channel
	userClient      UserClient
}

func NewNotificationService(
	emailService email.EmailService,
	storageService storage.MinIOStorage,
	rabbitURL string,
	userClient UserClient,
) (NotificationService, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	queues := []string{
		"booking_confirmations",
		"reservation_expired",
		"screening_missed",
		"ticket_emails",
	}

	for _, q := range queues {
		_, err := channel.QueueDeclare(
			q,     //name
			true,  //durable
			false, //delete when unused
			false, //exclusive
			false, //no-wait
			nil,   //arguments
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare queue %s: %w", q, err)
		}
	}

	ns := &notificationService{
		emailService:   emailService,
		storageService: storageService,
		rabbitConn:     conn,
		rabbitChannel:  channel,
		userClient:     userClient,
	}

	//background consumers
	go ns.consumeNotifications()

	return ns, nil
}

//sends confirmation for multiple tickets
func (s *notificationService) SendBookingConfirmation(ctx context.Context, userID int, tickets []domain.Ticket) error {
	message := BookingConfirmationMessage{
		UserID:  userID,
		Tickets: tickets,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return s.publishMessage(ctx, "booking_confirmations", body)
}

//sends notification when reservation expires
func (s *notificationService) SendReservationExpired(ctx context.Context, userID int, ticketID int) error {
	message := ReservationExpiredMessage{
		UserID:   userID,
		TicketID: ticketID,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return s.publishMessage(ctx, "reservation_expired", body)
}

//notification when user misses screening
func (s *notificationService) SendScreeningMissed(ctx context.Context, userID int, ticketID int) error {
	message := ScreeningMissedMessage{
		UserID:   userID,
		TicketID: ticketID,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return s.publishMessage(ctx, "screening_missed", body)
}

//sends individual ticket email with PDF
func (s *notificationService) SendTicketEmail(ctx context.Context, userEmail string, ticket *domain.TicketWithDetails, pdfData []byte, downloadURL string) error {
	return s.emailService.SendTicketEmail(ctx, userEmail, ticket, pdfData, downloadURL)
}

func (s *notificationService) publishMessage(ctx context.Context, queueName string, body []byte) error {
	err := s.rabbitChannel.PublishWithContext(
		ctx,
		"",        //exchange
		queueName, //routing key
		false,     //mandatory
		false,     //immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Persistent:  true,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

//notification messages from queues
func (s *notificationService) consumeNotifications() {
	go s.consumeQueue("booking_confirmations", s.processBookingConfirmation)

	go s.consumeQueue("reservation_expired", s.processReservationExpired)

	go s.consumeQueue("screening_missed", s.processScreeningMissed)
}

func (s *notificationService) consumeQueue(queueName string, handler func([]byte) error) {
	msgs, err := s.rabbitChannel.Consume(
		queueName,
		"",        //consumer
		false,     //auto-ack
		false,     //exclusive
		false,     //no-local
		false,     //no-wait
		nil,       //args
	)
	if err != nil {
		log.Printf("Failed to register consumer for %s: %v", queueName, err)
		return
	}

	for msg := range msgs {
		if err := handler(msg.Body); err != nil {
			log.Printf("Error processing message from %s: %v", queueName, err)
			msg.Nack(false, true) //requeue on error
		} else {
			msg.Ack(false)
		}
	}
}

func (s *notificationService) processBookingConfirmation(body []byte) error {
	var message BookingConfirmationMessage
	if err := json.Unmarshal(body, &message); err != nil {
		return err
	}

	user, err := s.userClient.GetUser(context.Background(), message.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	log.Printf("Processing booking confirmation for user %d: %d tickets", message.UserID, len(message.Tickets))


	_ = user
	return nil
}

func (s *notificationService) processReservationExpired(body []byte) error {
	var message ReservationExpiredMessage
	if err := json.Unmarshal(body, &message); err != nil {
		return err
	}

	user, err := s.userClient.GetUser(context.Background(), message.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	log.Printf("Sending reservation expired email to user %d for ticket %d", message.UserID, message.TicketID)

	return s.emailService.SendReservationExpired(context.Background(), user.Email, message.TicketID)
}

func (s *notificationService) processScreeningMissed(body []byte) error {
	var message ScreeningMissedMessage
	if err := json.Unmarshal(body, &message); err != nil {
		return err
	}

	user, err := s.userClient.GetUser(context.Background(), message.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	log.Printf("Sending screening missed email to user %d for ticket %d", message.UserID, message.TicketID)

	_ = user
	return nil
}

type BookingConfirmationMessage struct {
	UserID  int             `json:"user_id"`
	Tickets []domain.Ticket `json:"tickets"`
}

type ReservationExpiredMessage struct {
	UserID   int `json:"user_id"`
	TicketID int `json:"ticket_id"`
}

type ScreeningMissedMessage struct {
	UserID   int `json:"user_id"`
	TicketID int `json:"ticket_id"`
}

type UserClient interface {
	GetUser(ctx context.Context, userID int) (*User, error)
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *notificationService) Close() error {
	if s.rabbitChannel != nil {
		s.rabbitChannel.Close()
	}
	if s.rabbitConn != nil {
		s.rabbitConn.Close()
	}
	return nil
}