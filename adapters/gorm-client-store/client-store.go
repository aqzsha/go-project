package gorm_client_store

import (
	"auth/internal/app/models"
	"auth/pkg/graylog/logger"
	"context"
	"errors"
	"fmt"

	"github.com/go-oauth2/oauth2/v4"
	oauthErrs "github.com/go-oauth2/oauth2/v4/errors"
	oauthModels "github.com/go-oauth2/oauth2/v4/models"
	"gorm.io/gorm"
)

type GormClientStore struct {
	db *gorm.DB
}

func NewGormClientStore(db *gorm.DB) *GormClientStore {
	return &GormClientStore{db: db}
}

func (s *GormClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	var c models.OAuthClient
	if err := s.db.WithContext(ctx).First(&c, "id = ?", id).Error; err != nil {
		logger.Log.Error(fmt.Printf("failed GormClientStore method GetByID: %s", err.Error()))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, oauthErrs.ErrInvalidClient
		}

		return nil, err
	}
	return &oauthModels.Client{
		ID:     c.ID,
		Secret: c.Secret,
	}, nil
}
