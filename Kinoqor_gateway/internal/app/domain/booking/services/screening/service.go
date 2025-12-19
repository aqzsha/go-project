package screening

import (
	"context"
	microservice "gateway/internal/app/core/microservices"
	"io"
)

type Service interface {
	Create(ctx context.Context, payload io.Reader) microservice.Response
	Get(ctx context.Context, id string) microservice.Response
	Delete(ctx context.Context, id string) microservice.Response
	ListByFilm(ctx context.Context, id string) microservice.Response
	ListByCinema(ctx context.Context, id string) microservice.Response
	ListByHall(ctx context.Context, id string) microservice.Response
}

type service struct {
	microservice *microservice.RequestHandler
}

func NewService(microservice *microservice.RequestHandler) Service {
	return &service{microservice: microservice}
}

func (s *service) Create(ctx context.Context, payload io.Reader) microservice.Response {
	return s.microservice.Post(ctx, "screening/store", microservice.RequestOption{
		Body: payload,
	})
}

func (s *service) Get(ctx context.Context, id string) microservice.Response {
	return s.microservice.Get(ctx, "screening/get/"+id, microservice.RequestOption{})
}

func (s *service) Delete(ctx context.Context, id string) microservice.Response {
	return s.microservice.Delete(ctx, "screening/delete/"+id, microservice.RequestOption{})
}

func (s *service) ListByFilm(ctx context.Context, id string) microservice.Response {
	return s.microservice.Get(ctx, "screening/list-by-film/"+id, microservice.RequestOption{})
}

func (s *service) ListByCinema(ctx context.Context, id string) microservice.Response {
	return s.microservice.Get(ctx, "screening/list-by-cinema/"+id, microservice.RequestOption{})
}

func (s *service) ListByHall(ctx context.Context, id string) microservice.Response {
	return s.microservice.Get(ctx, "screening/list-by-hall/"+id, microservice.RequestOption{})
}
