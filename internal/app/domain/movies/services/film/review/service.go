package film

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
	ListFilm(ctx context.Context, id string) microservice.Response
}

type service struct {
	microservice *microservice.RequestHandler
}

func NewService(microservice *microservice.RequestHandler) Service {
	return &service{microservice: microservice}
}

func (s *service) Create(ctx context.Context, payload io.Reader) microservice.Response {
	return s.microservice.Post(ctx, "film/review/store", microservice.RequestOption{
		Body: payload,
	})
}

func (s *service) Get(ctx context.Context, id string) microservice.Response {
	return s.microservice.Get(ctx, "film/review/get/"+id, microservice.RequestOption{})
}

func (s *service) Delete(ctx context.Context, id string) microservice.Response {
	return s.microservice.Delete(ctx, "film/review/delete/"+id, microservice.RequestOption{})
}

func (s *service) List(ctx context.Context) microservice.Response {
	return s.microservice.Get(ctx, "film/review/list", microservice.RequestOption{})
}
func (s *service) ListFilm(ctx context.Context,id string) microservice.Response {
	return s.microservice.Get(ctx, "film/review/list/"+id, microservice.RequestOption{})
}
