package user

import (
	repoDto "auth/internal/app/domain/core/dto/repositories/user"
	dto "auth/internal/app/domain/core/dto/services/user"
	"auth/internal/app/domain/core/helpers/password"
	repository "auth/internal/app/domain/repositories/user"
	"auth/internal/app/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("user not found")
)

type Service interface {
	Create(ctx context.Context, input dto.CreateDTO) (models.User, error)
	Delete(ctx context.Context, input dto.DeleteDTO) (bool, error)
}

type service struct {
	db         *gorm.DB
	repository repository.Repository
}

func NewService(
	db *gorm.DB,
) Service {
	return &service{
		db:         db,
		repository: repository.NewRepository(db),
	}
}

func (s *service) Create(ctx context.Context, input dto.CreateDTO) (models.User, error) {
	passwd, err := password.Hash(input.Password)
	if err != nil{
		return models.User{}, err
	}
	return s.repository.Create(ctx, repoDto.CreateDTO{
		Email:    input.Email,
		Password: passwd,
	})
}

func (s *service) Delete(ctx context.Context, input dto.DeleteDTO) (bool, error) {
	ok, err := s.repository.Delete(ctx, repoDto.DeleteDTO{
		Email: input.Email,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ok, ErrNotFound
		}
	}

	return ok, err
}
