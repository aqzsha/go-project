package film

import (
	"context"
	microservice "gateway/internal/app/core/microservices"
	"io"
)

type Service interface {
	Create(ctx context.Context, payload io.Reader) microservice.Response
	Get(ctx context.Context, id string) microservice.Response
	List(ctx context.Context) microservice.Response
	Delete(ctx context.Context, id string) microservice.Response
}

type service struct {
	requestHandler *microservice.RequestHandler
}

func NewService(rh *microservice.RequestHandler) Service {
	return &service{requestHandler: rh}
}

func (s *service) Create(ctx context.Context, payload io.Reader) microservice.Response {
	return s.requestHandler.Post(ctx, "film/store", microservice.RequestOption{
		Body: payload,
	})
}

func (s *service) Get(ctx context.Context, id string) microservice.Response {
	return s.requestHandler.Get(ctx, "film/get/"+id, microservice.RequestOption{})
}

func (s *service) List(ctx context.Context) microservice.Response {
	return s.requestHandler.Get(ctx, "film/list", microservice.RequestOption{})
}

func (s *service) Delete(ctx context.Context, id string) microservice.Response {
	return s.requestHandler.Delete(ctx, "film/delete/"+id, microservice.RequestOption{})
}
