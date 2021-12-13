package transport

import (
	"context"

	"github.com/c95rt/bootcamp-user/grpc/endpoints"
	"github.com/c95rt/bootcamp-user/grpc/pb"
	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type gRPCServer struct {
	auth       gt.Handler
	createUser gt.Handler
	pb.UnimplementedUserServiceServer
}

func NewGRPCServer(endpoints endpoints.Endpoints, logger log.Logger) pb.UserServiceServer {
	return &gRPCServer{
		auth: gt.NewServer(
			endpoints.Authenticate,
			decodeAuthRequest,
			encodeAuthResponse,
		),
		createUser: gt.NewServer(
			endpoints.CreateUser,
			decodeCreateUserRequest,
			encodeCreateUserResponse,
		),
	}
}

func (s *gRPCServer) Authenticate(ctx context.Context, req *pb.AuthRequest) (*pb.AuthReply, error) {
	_, resp, err := s.auth.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.AuthReply), nil
}

func decodeAuthRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.AuthRequest)
	return endpoints.AuthReq{Pwd: req.PwdHash, Name: req.Name}, nil
}

func encodeAuthResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*endpoints.AuthRes)
	return &pb.AuthReply{ID: resp.ID}, nil
}

func (s *gRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateUserResponse), nil
}

func decodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.CreateUserRequest)
	return endpoints.CreateUserReq{
		Name:    req.Name,
		Pwd:     req.Pwd,
		Age:     req.Age,
		AddInfo: req.AddInfo,
	}, nil
}

func encodeCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*endpoints.CreateUserResponse)
	return &pb.CreateUserResponse{ID: resp.UserId}, nil
}
