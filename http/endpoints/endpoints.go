package endpoints

import (
	"context"

	"github.com/c95rt/bootcamp-user/http/models"
	"github.com/c95rt/bootcamp-user/http/service"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Login      endpoint.Endpoint
	InsertUser endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Login:      makeLoginEndpoint(s),
		InsertUser: makeInsertUserEndpoint(s),
	}
}

func makeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.LoginRequest)
		user, err := s.Login(ctx, &req)

		return user, err
	}
}

func makeInsertUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.InsertUserRequest)
		user, err := s.InsertUser(ctx, &req)

		return user, err
	}
}
