package cinema

import (
	"context"
	microservice "gateway/internal/app/core/microservices"
	"io"
)

type Service interface {
	Create(ctx context.Context, payload io.Reader) microservice.Response
	Get(ctx context.Context, id string) microservice.Response
	Delete(ctx context.Context, id string) microservice.Response
	List(ctx context.Context) microservice.Response	
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
		"cinema/store",
		microservice.RequestOption{
			Body: payload,
		},
	)
}

func (s *service) Get(ctx context.Context, id string) microservice.Response {
	return s.microservice.Get(
		ctx,
		"cinema/get/"+id,
		microservice.RequestOption{},
	)
}

func (s *service) Delete(ctx context.Context, id string) microservice.Response {
	return s.microservice.Delete(
		ctx,
		"cinema/delete/"+id,
		microservice.RequestOption{},
	)
}

func (s *service) List(ctx context.Context) microservice.Response {
	return s.microservice.Get(ctx, "cinema/list", microservice.RequestOption{})
}