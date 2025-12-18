package microservices

import (
	"fmt"
	"movies/configs"
	myHttp "movies/internal/app/core/http"
	"movies/internal/app/core/microservices/booking"
	"movies/internal/app/core/retry"
	"time"
)

type Clients struct {
	Booking *booking.HTTPClient
}

func NewClients() (*Clients, error) {
	base := myHttp.New(myHttp.ClientConfig{
		Timeout: time.Duration(configs.Config.Microservices.Http.Timeout) * time.Second,
	})

	r := &retry.Retrier{
		Base:    base,
		Retries: configs.Config.Microservices.Http.Retries,
		Backoff: time.Duration(configs.Config.Microservices.Http.BackoffMillis) * time.Millisecond,
	}

	bookingCli, err := booking.New(configs.Config.Microservices.Booking.BaseURL, r)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking client: %w", err)
	}

	return &Clients{Booking: bookingCli}, nil
}
