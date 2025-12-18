package genre

import (
	"context"
	microservice "gateway/internal/app/core/microservices"
	"io"
)

type Service interface {
	Create(ctx context.Context, payload io.Reader) microservice.Response
	Get(ctx context.Context, id string) microservice.Response
	Delete(ctx context.Context, id string) microservice.Response
}

type service struct {
	microservice *microservice.RequestHandler
}

func NewService(microservice *microservice.RequestHandler) Service {
	return &service{microservice: microservice}
}

func (s *service) Create(ctx context.Context, payload io.Reader) microservice.Response {
	return s.microservice.Post(
		ctx,
		"genre/store",
		microservice.RequestOption{
			Body: payload,
		},
	)
}

func (s *service) Get(ctx context.Context, id string) microservice.Response {
	return s.microservice.Get(
		ctx,
		"genre/get/"+id,
		microservice.RequestOption{},
	)
}

func (s *service) Delete(ctx context.Context, id string) microservice.Response {
	return s.microservice.Delete(
		ctx,
		"genre/delete/"+id,
		microservice.RequestOption{},
	)
}
