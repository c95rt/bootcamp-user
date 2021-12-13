package endpoints

import (
	"context"

	"github.com/c95rt/bootcamp-user/grpc/repository"
	"github.com/c95rt/bootcamp-user/grpc/service"
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
		req := request.(repository.LoginRequest)
		return s.Login(ctx, &req)
	}
}

func makeInsertUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(repository.InsertUserRequest)
		return s.InsertUser(ctx, &req)
	}
}
