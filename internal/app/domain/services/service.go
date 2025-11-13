package services

import (
	"auth/configs"
	dto "auth/internal/app/domain/core/dto/services"
	repository "auth/internal/app/domain/repositories"
	tokenService "auth/internal/app/domain/services/token"
	"auth/internal/app/models"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-oauth2/oauth2/v4"
	oauthErrs "github.com/go-oauth2/oauth2/v4/errors"
	"gorm.io/gorm"
)

var (
	ErrNotFound        = errors.New("user not found")
	ErrFailedLogin     = errors.New("login failed")
	ErrUnauthenticated = errors.New("unauthenticated")
)

type Service interface {
	Login(ctx context.Context, dto dto.LoginDTO) (tokenService.Response, int, error)
	RefreshToken(ctx context.Context, refresh string) (tokenService.Response, error)
	CheckToken(ctx context.Context, access string) (models.User, error)
	Logout(ctx context.Context, access string) (bool, error)
}

type service struct {
	db           *gorm.DB
	repository   repository.Repository
	tokenService tokenService.Service
	tokenStore   oauth2.TokenStore
}

func NewService(db *gorm.DB, tokenService tokenService.Service, tokenStore oauth2.TokenStore) Service {
	return &service{
		db:           db,
		repository:   repository.NewRepository(db),
		tokenService: tokenService,
		tokenStore:   tokenStore,
	}
}

func (s *service) Login(ctx context.Context, input dto.LoginDTO) (tokenService.Response, int, error) {
	token, err := s.tokenService.IssueToken(ctx, map[string]string{
		"grant_type":    "password",
		"client_id":     configs.Config.Oauth.ClientID,
		"client_secret": configs.Config.Oauth.ClientSecret,
		"username":      input.Email,
		"password":      input.Password,
	})
	if err != nil {
		return tokenService.Response{}, http.StatusNotFound, ErrFailedLogin
	}

	return token, http.StatusOK, err
}

func (s *service) RefreshToken(ctx context.Context, refresh string) (tokenService.Response, error) {
	token, err := s.tokenService.IssueToken(ctx, map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     configs.Config.Oauth.ClientID,
		"client_secret": configs.Config.Oauth.ClientSecret,
		"refresh_token": refresh,
	})
	if err != nil {
		if strings.Contains(err.Error(), oauthErrs.ErrInvalidGrant.Error()) {
			return tokenService.Response{}, ErrUnauthenticated
		}
	}

	return token, err
}

func (s *service) CheckToken(ctx context.Context, access string) (models.User, error) {
	tokenInfo, err := s.tokenService.ValidateToken(ctx, access)
	if err != nil {
		return models.User{}, ErrUnauthenticated
	}

	intUserID, err := strconv.ParseInt(tokenInfo.GetUserID(), 10, 64)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to convert user id to int: %w", err)
	}

	return s.repository.GetById(ctx, intUserID)
}

func (s *service) Logout(ctx context.Context, access string) (bool, error) {
	tokenInfo, err := s.tokenService.ValidateToken(ctx, access)
	if err != nil {
		return false, ErrUnauthenticated
	}

	if err = s.tokenStore.RemoveByAccess(ctx, tokenInfo.GetAccess()); err != nil {
		return false, fmt.Errorf("failed to remove access: %w", err)
	}

	return true, nil
}
