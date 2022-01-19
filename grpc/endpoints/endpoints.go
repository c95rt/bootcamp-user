package endpoints

import (
	"context"

	"github.com/c95rt/bootcamp-user/grpc/models"
	"github.com/c95rt/bootcamp-user/grpc/service"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Login      endpoint.Endpoint
	InsertUser endpoint.Endpoint
	UpdateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
	DeleteUser endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Login:      makeLoginEndpoint(s),
		InsertUser: makeInsertUserEndpoint(s),
		UpdateUser: makeUpdateUserEndpoint(s),
		GetUser:    makeGetUserEndpoint(s),
		DeleteUser: makeDeleteUserEndpoint(s),
	}
}

func makeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.LoginRequest)
		return s.Login(ctx, &req)
	}
}

func makeInsertUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.InsertUserRequest)
		return s.InsertUser(ctx, &req)
	}
}

func makeUpdateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.UpdateUserRequest)
		return s.UpdateUser(ctx, &req)
	}
}

func makeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.GetUserRequest)
		return s.GetUser(ctx, &req)
	}
}

func makeDeleteUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.DeleteUserRequest)
		return s.DeleteUser(ctx, &req)
	}
}
