package password

import (
	"auth/internal/app/models"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var (
	ErrNotFound        = errors.New("token not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrUnauthenticated = errors.New("unauthenticated")
)

type Repository interface {
	UpdatePassword(ctx context.Context, userID int64, hashed string) error
	GetResetToken(ctx context.Context, token string, email string) (models.User, error)
	GetById(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	CreateResetToken(ctx context.Context, email string, pinCode string, token string) error
	DeleteResetToken(ctx context.Context, token string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateResetToken(ctx context.Context, email string, pinCode string, token string) error {
	reset := models.PasswordReset{
		Email:   email,
		PinCode: pinCode,
		Token:   token,
	}

	err := r.db.WithContext(ctx).Create(&reset).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get create reset token: %w", err)
	}

	return nil
}

func (r *repository) GetResetToken(ctx context.Context, token string, email string) (models.User, error) {
	var reset models.PasswordReset
	if err := r.db.WithContext(ctx).Where("token = ? AND email = ?", token, email).First(&reset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrInvalidToken
		}
		return models.User{}, fmt.Errorf("failed to get reset password: %w", err)
	}

	if time.Since(reset.CreatedAt) > time.Hour {
		return models.User{}, ErrTokenExpired
	}

	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, ErrUserNotFound
	}

	return user, nil
}

func (r *repository) DeleteResetToken(ctx context.Context, token string) error {
	err := r.db.WithContext(ctx).Where("token = ?", token).Delete(&models.PasswordReset{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete reset token: %w", err)
	}
	return nil
}

func (r *repository) UpdatePassword(ctx context.Context, userID int64, hashed string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password", hashed).Error
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

func (r *repository) GetById(ctx context.Context, id int64) (*models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get User by id: %w", err)
	}

	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get User by email: %w", err)
	}

	return &user, nil
}
