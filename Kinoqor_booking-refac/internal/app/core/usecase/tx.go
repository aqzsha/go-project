package usecase

import "context"

// TxManager запускает fn внутри БД-транзакции.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
