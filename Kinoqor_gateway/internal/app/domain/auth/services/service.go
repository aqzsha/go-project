package services

import (
	"context"
	"fmt"
	"gateway/internal/app/core/contracts/headercontract"
	microservice "gateway/internal/app/core/microservices"
	"gateway/pkg/parsertructure"
	"io"
	"net/http"
)

type Service interface {
	Login(ctx context.Context, payload io.Reader) microservice.Response
	Refresh(ctx context.Context, payload io.Reader) microservice.Response
	Check(ctx context.Context) microservice.Response
	Logout(ctx context.Context) microservice.Response
	GetAuthorizedID(ctx context.Context) (int64, error)
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

func (s *service) GetAuthorizedID(ctx context.Context) (int64, error) {
	type responseBody struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}

	resp := s.Check(ctx)
	if resp.Error != nil || resp.StatusCode >= http.StatusBadRequest {
		return 0, fmt.Errorf("failed to get email: %w", resp.Error)
	}

	var response responseBody

	if err := parsertructure.Parse[responseBody](resp.Data, &response); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return response.Data.ID, nil
}

