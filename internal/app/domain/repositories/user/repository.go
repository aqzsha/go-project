package user

import (
	dto "auth/internal/app/domain/core/dto/repositories/user"
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
	Create(ctx context.Context, input dto.CreateDTO) (models.User, error)
	Delete(ctx context.Context, input dto.DeleteDTO) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input dto.CreateDTO) (models.User, error) {
	user := models.User{
		Email:    input.Email,
		Password: input.Password,
	}

	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("failed to create User: %w", err)
	}

	return user, nil
}

func (r *repository) Delete(ctx context.Context, input dto.DeleteDTO) (bool, error) {
	res := r.db.WithContext(ctx).Where("email = ?", input.Email).Delete(&models.User{})
	if res.Error != nil {
		return false, fmt.Errorf("failed to delete User: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return false, ErrNotFound
	}

	return true, nil
}
