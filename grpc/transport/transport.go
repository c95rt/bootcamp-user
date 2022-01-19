package transport

import (
	"context"
	"strconv"

	"github.com/c95rt/bootcamp-user/grpc/endpoints"
	"github.com/c95rt/bootcamp-user/grpc/models"
	"github.com/c95rt/bootcamp-user/grpc/pb"
	gt "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	login      gt.Handler
	insertUser gt.Handler
	updateUser gt.Handler
	getUser    gt.Handler
	deleteUser gt.Handler
	pb.UnimplementedUserServiceServer
}

func NewGRPCServer(endpoints endpoints.Endpoints) pb.UserServiceServer {
	return &gRPCServer{
		login: gt.NewServer(
			endpoints.Login,
			decodeLoginRequest,
			encodeLoginResponse,
		),
		insertUser: gt.NewServer(
			endpoints.InsertUser,
			decodeInsertUserRequest,
			encodeInsertUserResponse,
		),
		updateUser: gt.NewServer(
			endpoints.UpdateUser,
			decodeUpdateUserRequest,
			encodeUpdateUserResponse,
		),
		getUser: gt.NewServer(
			endpoints.GetUser,
			decodeGetUserRequest,
			encodeGetUserResponse,
		),
		deleteUser: gt.NewServer(
			endpoints.DeleteUser,
			decodeDeleteUserRequest,
			encodeDeleteUserResponse,
		),
	}
}

func (s *gRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, resp, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LoginResponse), nil
}

func decodeLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.LoginRequest)
	return models.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func encodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*models.LoginResponse)
	return &pb.LoginResponse{
		Token: resp.Token,
	}, nil
}

func (s *gRPCServer) InsertUser(ctx context.Context, req *pb.InsertUserRequest) (*pb.InsertUserResponse, error) {
	_, resp, err := s.insertUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.InsertUserResponse), nil
}

func decodeInsertUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.InsertUserRequest)
	return models.InsertUserRequest{
		Email:     req.Email,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Password:  req.Password,
		BirthDate: req.BirthDate,
		Address:   req.Address,
	}, nil
}

func encodeInsertUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*models.InsertUserResponse)
	return &pb.InsertUserResponse{
		Id:        strconv.Itoa(resp.ID),
		Email:     resp.Email,
		Firstname: resp.Firstname,
		Lastname:  resp.Lastname,
		Password:  resp.Password,
		BirthDate: resp.BirthDate,
		Address:   resp.Address,
	}, nil
}

func (s *gRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	_, resp, err := s.updateUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UpdateUserResponse), nil
}

func decodeUpdateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UpdateUserRequest)
	userID, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}
	return models.UpdateUserRequest{
		ID:        userID,
		Email:     req.Email,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Password:  req.Password,
		BirthDate: req.BirthDate,
		Address:   req.Address,
	}, nil
}

func encodeUpdateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*models.UpdateUserResponse)
	return &pb.UpdateUserResponse{
		Id:        strconv.Itoa(resp.ID),
		Email:     resp.Email,
		Firstname: resp.Firstname,
		Lastname:  resp.Lastname,
		Password:  resp.Password,
		BirthDate: resp.BirthDate,
		Address:   resp.Address,
	}, nil
}

func (s *gRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	_, resp, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetUserResponse), nil
}

func decodeGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetUserRequest)
	userID, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}
	return models.GetUserRequest{
		ID: userID,
	}, nil
}

func encodeGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*models.GetUserResponse)
	return &pb.GetUserResponse{
		Id:        strconv.Itoa(resp.ID),
		Email:     resp.Email,
		Firstname: resp.Firstname,
		Lastname:  resp.Lastname,
		Password:  resp.Password,
		BirthDate: resp.BirthDate,
		Address:   resp.Address,
	}, nil
}

func (s *gRPCServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, resp, err := s.deleteUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.DeleteUserResponse), nil
}

func decodeDeleteUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.DeleteUserRequest)
	userID, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}
	return models.DeleteUserRequest{
		ID: userID,
	}, nil
}

func encodeDeleteUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*models.DeleteUserResponse)
	return &pb.DeleteUserResponse{
		Id: strconv.Itoa(resp.ID),
	}, nil
}
