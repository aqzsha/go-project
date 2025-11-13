package password

import (
	dto "auth/internal/app/domain/core/dto/services/password"
	"auth/internal/app/domain/core/helpers/mailer"
	"auth/internal/app/domain/core/helpers/token"
	"auth/internal/app/domain/core/helpers/validate"
	repository "auth/internal/app/domain/repositories/password"
	tokenService "auth/internal/app/domain/services/token"
	"auth/internal/app/models"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrNotFound        = errors.New("token not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrInvalidUser     = errors.New("invalid user id")
	ErrInvalidPassword = errors.New("invalid password")
)

type Service interface {
	ForgotPassword(ctx context.Context, input dto.ForgotPasswordDTO) error
	ResetPassword(ctx context.Context, input dto.ResetPasswordDTO) error
	ChangePassword(ctx context.Context, input dto.ChangePasswordDTO) error
	GetById(ctx context.Context, userID int64) (*models.User, error)
}

type service struct {
	db           *gorm.DB
	repository   repository.Repository
	tokenService tokenService.Service
}

func NewService(db *gorm.DB, tokenService tokenService.Service) Service {
	return &service{
		db:           db,
		repository:   repository.NewRepository(db),
		tokenService: tokenService,
	}
}

func (s *service) ForgotPassword(ctx context.Context, input dto.ForgotPasswordDTO) error {
	user, err := s.getByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUnauthenticated) {
			return ErrUnauthenticated
		}
		return err
	}

	randomToken := token.GenerateRandomToken(32)
	expiration := time.Now().Add(1 * time.Hour)

	if err = s.repository.CreateResetToken(ctx, user.Email, randomToken, expiration); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	resetLink := fmt.Sprintf("https://yourfrontend.app/reset-password?token=%s", randomToken)
	email := mailer.Mailer{
		"armankaliakyn@gmail.com",
		"arman2002",
		"localhost",
		"80",
	}
	err = email.SendMail(input.Email, "Password Reset", fmt.Sprintf("Click here to reset your password: %s", resetLink))
	if err != nil {
		return fmt.Errorf("failed to send email")
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, input dto.ResetPasswordDTO) error {
	user, reset, err := s.repository.GetByResetToken(ctx, input.Token)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		if errors.Is(err, repository.ErrInvalidToken) {
			return ErrInvalidToken
		}
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := s.repository.UpdatePassword(ctx, user.ID, string(hashed)); err != nil {
		return err
	}

	s.repository.UpdateTokenUsed(ctx, reset.Token)
	return err
}

func (s *service) ChangePassword(ctx context.Context, input dto.ChangePasswordDTO) error {
	tokenInfo, err := s.tokenService.ValidateToken(ctx, input.Token)
	if err != nil {
		return ErrInvalidUser
	}

	u64, err := strconv.ParseUint(tokenInfo.GetUserID(), 10, 64)
	if err != nil {
		return ErrInvalidUser
	}

	userID := int64(u64)
	user, err := s.GetById(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	err = validate.Validate(input, user)
	if err != nil {
		return ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	err = s.repository.UpdatePassword(ctx, user.ID, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetById(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.repository.GetById(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUnauthenticated
		}
	}
	return user, nil
}

func (s *service) getByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUnauthenticated) {
			return nil, ErrUnauthenticated
		}
		return nil, err
	}
	return user, nil
}
