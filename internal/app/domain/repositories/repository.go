package repositories

import (
	"auth/internal/app/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("user not found")
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetById(ctx context.Context, id int64) (models.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrNotFound
		}

		return models.User{}, fmt.Errorf("failed to get User by email: %w", err)
	}

	return user, nil
}

func (r *repository) GetById(ctx context.Context, id int64) (models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrNotFound
		}

		return models.User{}, fmt.Errorf("failed to get User by id: %w", err)
	}

	return user, nil
}
