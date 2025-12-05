package services

import (
	"context"
	"fmt"
	"gateway/internal/app/core/contracts/headercontract"
	microservice "gateway/internal/app/core/microservices"
	dto "gateway/internal/app/domain/auth/dto/auth"
	"gateway/pkg/parsertructure"
	"io"
	"net/http"
	"net/url"
)

type Service interface {
	Login(ctx context.Context, payload io.Reader) microservice.Response
	Refresh(ctx context.Context, payload io.Reader) microservice.Response
	Check(ctx context.Context) microservice.Response
	Logout(ctx context.Context) microservice.Response
	GetAuthorizedEmail(ctx context.Context) (string, error)
	GetByEmail(ctx context.Context, email string) (dto.AuthUser, error)
}

type service struct {
	requestHandler *microservice.RequestHandler
}

func NewService(requestHandler *microservice.RequestHandler) Service {
	return &service{requestHandler: requestHandler}
}

func (s *service) Login(ctx context.Context, payload io.Reader) microservice.Response {
	return s.requestHandler.Post(ctx, "login", microservice.RequestOption{
		Body: payload,
	})
}

func (s *service) Refresh(ctx context.Context, payload io.Reader) microservice.Response {
	return s.requestHandler.Post(ctx, "refresh", microservice.RequestOption{
		Body: payload,
	})
}

func (s *service) Logout(ctx context.Context) microservice.Response {
	return s.requestHandler.Post(ctx, "logout", microservice.RequestOption{
		Header: map[string]string{
			"Authorization": headercontract.GetBearer(ctx),
		},
	})
}

func (s *service) Check(ctx context.Context) microservice.Response {
	return s.requestHandler.Post(ctx, "check", microservice.RequestOption{
		Header: map[string]string{
			"Authorization": headercontract.GetBearer(ctx),
		},
	})
}

func (s *service) GetAuthorizedEmail(ctx context.Context) (string, error) {
	type responseBody struct {
		Data struct {
			Email string `json:"email"`
		} `json:"data"`
	}

	resp := s.Check(ctx)
	if resp.Error != nil || resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("failed to get email: %w", resp.Error)
	}

	var response responseBody

	if err := parsertructure.Parse[responseBody](resp.Data, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return response.Data.Email, nil
}


func (s *service) GetByEmail(ctx context.Context, email string) (dto.AuthUser, error) {
	resp := s.requestHandler.Get(ctx, "by-email", microservice.RequestOption{
		Query: url.Values{
			"email": []string{email},
		},
	})
	if resp.Error != nil || resp.StatusCode >= http.StatusBadRequest {
		return dto.AuthUser{}, fmt.Errorf("get user by email failed: %s", resp.Error)
	}

	type responseBody struct {
		Data struct {
			dto.AuthUser
		} `json:"data"`
	}

	var respBody responseBody
	if err := parsertructure.Parse[responseBody](resp.Data, &respBody); err != nil {
		return dto.AuthUser{}, fmt.Errorf("failed to parse auth data: %w", err)
	}

	return dto.AuthUser{
		ID:         respBody.Data.AuthUser.ID,
	}, nil
}