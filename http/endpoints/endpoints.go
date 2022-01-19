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
	UpdateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
	DeleteUser endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Login:      makeLoginEndpoint(s),
		InsertUser: makeInsertUserEndpoint(s),
		GetUser:    makeGetUserEndpoint(s),
		UpdateUser: makeUpdateUserEndpoint(s),
		DeleteUser: makeDeleteUserEndpoint(s),
	}
}

func makeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.LoginRequest)
		response, err = s.Login(ctx, &req)

		return response, err
	}
}

func makeInsertUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.InsertUserRequest)
		response, err = s.InsertUser(ctx, &req)

		return response, err
	}
}

func makeUpdateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.UpdateUserRequest)
		response, err = s.UpdateUser(ctx, &req)

		return response, err
	}
}

func makeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.GetUserRequest)
		response, err = s.GetUser(ctx, &req)

		return response, err
	}
}

func makeDeleteUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.DeleteUserRequest)
		response, err = s.DeleteUser(ctx, &req)

		return response, err
	}
}
