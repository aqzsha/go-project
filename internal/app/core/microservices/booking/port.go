package booking

import (
	"context"
)

type Client interface {
	StoreSeat(ctx context.Context, input StoreSeat) error
	DeleteSeat(ctx context.Context, id int64) error
}
