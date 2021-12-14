package transport

import (
	"context"
	"strconv"

	"github.com/c95rt/bootcamp-user/grpc/endpoints"
	"github.com/c95rt/bootcamp-user/grpc/pb"
	"github.com/c95rt/bootcamp-user/grpc/repository"
	gt "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	login      gt.Handler
	insertUser gt.Handler
	pb.UnimplementedUserServiceServer
}

func NewGRPCServer(endpoints endpoints.Endpoints) pb.UserServiceServer {
	return &gRPCServer{
		login: gt.NewServer(
			endpoints.Login,
			decodeAuthRequest,
			encodeAuthResponse,
		),
		insertUser: gt.NewServer(
			endpoints.InsertUser,
			decodeInsertUserRequest,
			encodeInsertUserResponse,
		),
	}
}

func (s *gRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.User, error) {
	_, resp, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.User), nil
}

func decodeAuthRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.LoginRequest)
	return repository.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func encodeAuthResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*repository.User)
	return &pb.User{
		ID: strconv.Itoa(resp.ID),
	}, nil
}

func (s *gRPCServer) InsertUser(ctx context.Context, req *pb.InsertUserRequest) (*pb.User, error) {
	_, resp, err := s.insertUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.User), nil
}

func decodeInsertUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.InsertUserRequest)
	return repository.InsertUserRequest{
		Email:     req.Email,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Password:  req.Password,
		BirthDate: req.BirthDate,
		Address:   req.Address,
	}, nil
}

func encodeInsertUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*repository.User)
	return &pb.User{
		ID: strconv.Itoa(resp.ID),
	}, nil
}
