package password

import (
	send_email "auth/internal/app/core/helpers/send-email"
	dto "auth/internal/app/domain/core/dto/services/password"
	emailDTO "auth/internal/app/domain/core/dto/services/send-email"
	"auth/internal/app/domain/core/helpers/token"
	"auth/internal/app/domain/core/helpers/validate"
	repository "auth/internal/app/domain/repositories/password"
	tokenService "auth/internal/app/domain/services/token"
	"auth/internal/app/models"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"

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

	pinCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	randomToken := token.GenerateRandomToken(32)
	data := emailDTO.EmailData{
		Subject: "Сброс пароля",
		PinCode: pinCode,
	}

	if err = s.repository.CreateResetToken(ctx, user.Email, pinCode, randomToken); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	err = send_email.SendEmail(emailDTO.EmailDTO{
		Email:    user.Email,
		Template: "reset-password",
		Data:     data,
	})

	return err
}

func (s *service) ResetPassword(ctx context.Context, input dto.ResetPasswordDTO) error {
	user, err := s.repository.GetResetToken(ctx, input.Token, input.Email)
	if err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := s.repository.UpdatePassword(ctx, user.ID, string(hashed)); err != nil {
		return err
	}

	if err := s.repository.DeleteResetToken(ctx, input.Token); err != nil {
		return err
	}

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
