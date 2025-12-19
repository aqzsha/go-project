package book

import (
	"context"
	microservice "gateway/internal/app/core/microservices"
	"io"
)

type Service interface {
	Create(ctx context.Context, payload io.Reader) microservice.Response
	Delete(ctx context.Context, id string) microservice.Response
	Get(ctx context.Context, id string) microservice.Response
	PaidTicket(ctx context.Context, id string) microservice.Response
	ScanTicket(ctx context.Context, token string) microservice.Response
}

type service struct {
	microservice *microservice.RequestHandler
}

func NewService(microservice *microservice.RequestHandler) Service {
	return &service{microservice: microservice}
}

func (s *service) Create(ctx context.Context, payload io.Reader) microservice.Response {
	return s.microservice.Post(ctx, "book/store", microservice.RequestOption{
		Body: payload,
	})
}

func (s *service) Delete(ctx context.Context, id string) microservice.Response {
	return s.microservice.Delete(ctx, "book/delete/"+id, microservice.RequestOption{})
}

func (s *service) Get(ctx context.Context, id string) microservice.Response {
	return s.microservice.Get(ctx, "book/get/"+id, microservice.RequestOption{})
}

func (s *service) PaidTicket(ctx context.Context, id string) microservice.Response {
	return s.microservice.Post(ctx, "book/paid-ticket/"+id, microservice.RequestOption{})
}

func (s *service) ScanTicket(ctx context.Context, token string) microservice.Response {
	return s.microservice.Get(ctx, "book/tickets/scan/"+token, microservice.RequestOption{})
}
