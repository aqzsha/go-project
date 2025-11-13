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
	CreateResetToken(ctx context.Context, email, token string, expiration time.Time) error
	UpdatePassword(ctx context.Context, userID int64, hashed string) error
	GetByResetToken(ctx context.Context, token string) (*models.User, *models.PasswordReset, error)
	UpdateTokenUsed(ctx context.Context, token string)
	GetById(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateResetToken(ctx context.Context, email, token string, expiration time.Time) error {
	reset := models.PasswordReset{
		Email:     email,
		Token:     token,
		ExpiresAt: expiration,
		Used:      false,
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

func (r *repository) GetByResetToken(ctx context.Context, token string) (*models.User, *models.PasswordReset, error) {
	var reset models.PasswordReset
	if err := r.db.WithContext(ctx).Where("token = ? AND used = ?", token, false).First(&reset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrInvalidToken
		}
		return nil, nil, fmt.Errorf("failed to get reset password: %w", err)
	}

	if reset.ExpiresAt.Before(time.Now()) {
		return nil, nil, ErrTokenExpired
	}

	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", reset.Email).First(&user).Error; err != nil {
		return nil, nil, ErrUserNotFound
	}

	return &user, &reset, nil
}

func (r *repository) UpdatePassword(ctx context.Context, userID int64, hashed string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password", hashed).Error
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

func (r *repository) UpdateTokenUsed(ctx context.Context, token string) {
	r.db.WithContext(ctx).Model(&models.PasswordReset{}).Where("token = ?", token).Update("used", true)
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
