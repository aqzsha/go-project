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
	CreateGenre(ctx context.Context, payload io.Reader) microservice.Response
	DeleteGenre(ctx context.Context, id string) microservice.Response
	GetGenreList(ctx context.Context) microservice.Response
}

type service struct {
	microservice *microservice.RequestHandler
}

func NewService(microservice *microservice.RequestHandler) Service {
	return &service{microservice: microservice}
}

func (s *service) Create(ctx context.Context, payload io.Reader) microservice.Response {
	return s.microservice.Post(ctx, "film/store", microservice.RequestOption{
		Body: payload,
	})
}

func (s *service) Get(ctx context.Context, id string) microservice.Response {
	return s.microservice.Get(ctx, "film/get/"+id, microservice.RequestOption{})
}

func (s *service) Delete(ctx context.Context, id string) microservice.Response {
	return s.microservice.Delete(ctx, "film/delete/"+id, microservice.RequestOption{})
}

func (s *service) List(ctx context.Context) microservice.Response {
	return s.microservice.Get(ctx, "film/list", microservice.RequestOption{})
}

func (s *service) CreateGenre(ctx context.Context, payload io.Reader) microservice.Response {
	return s.microservice.Post(ctx, "film/genre/store", microservice.RequestOption{
		Body: payload,
	})
}

func (s *service) DeleteGenre(ctx context.Context, id string) microservice.Response {
	return s.microservice.Delete(ctx, "film/genre/delete/"+id, microservice.RequestOption{})
}

func (s *service) GetGenreList(ctx context.Context) microservice.Response {
	return s.microservice.Get(ctx, "film/genre/list", microservice.RequestOption{})
}

