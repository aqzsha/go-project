package password

import (
	"context"
	"gateway/internal/app/core/contracts/headercontract"
	microservice "gateway/internal/app/core/microservices"
	"io"
)

type Service interface {
	ForgotPassword(ctx context.Context, payload io.Reader) microservice.Response
	ResetPassword(ctx context.Context, payload io.Reader) microservice.Response
	ChangePassword(ctx context.Context, payload io.Reader) microservice.Response
	VerifyPin(ctx context.Context, payload io.Reader) microservice.Response
}

type service struct {
	requestHandler *microservice.RequestHandler
}

func NewService(requestHandler *microservice.RequestHandler) Service {
	return &service{requestHandler: requestHandler}
}

func (s *service) ForgotPassword(ctx context.Context, payload io.Reader) microservice.Response {
	return s.requestHandler.Post(ctx, "password/forgot", microservice.RequestOption{
		Body: payload,
	})
}

func (s *service) ResetPassword(ctx context.Context, payload io.Reader) microservice.Response {
	return s.requestHandler.Post(ctx, "password/reset", microservice.RequestOption{
		Body: payload,
	})
}

func (s *service) ChangePassword(ctx context.Context, payload io.Reader) microservice.Response {
	return s.requestHandler.Post(ctx, "password/change", microservice.RequestOption{
		Body: payload,
		Header: map[string]string{
			"Authorization": headercontract.GetBearer(ctx),
		},
	})
}

func (s *service) VerifyPin(ctx context.Context, payload io.Reader) microservice.Response {
	return s.requestHandler.Post(ctx, "password/verify-pin", microservice.RequestOption{
		Body: payload,
	})
}
