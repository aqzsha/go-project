package gormtx

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

func Inject(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func From(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(txKey{}).(*gorm.DB)

	return tx
}

type Manager struct {
	DB *gorm.DB
}

func NewManager(db *gorm.DB) *Manager { return &Manager{DB: db} }

func (m *Manager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxTx := Inject(ctx, tx)

		return fn(ctxTx)
	})
}
