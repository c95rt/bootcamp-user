package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/c95rt/bootcamp-user/grpc/pb"
	"github.com/c95rt/bootcamp-user/http/models"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type UserRepository interface {
	Login(ctx context.Context, request *models.LoginRequest) (*models.User, error)
	InsertUser(ctx context.Context, req *models.InsertUserRequest) (*models.User, error)
}

func (c *Conn) Login(ctx context.Context, request *models.LoginRequest) (*models.User, error) {
	logger := log.With(c.logger, "method", "Authenticate")

	client := pb.NewUserServiceClient(c.conn)

	user, err := client.Login(ctx, &pb.LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	userID, err := strconv.Atoi(user.ID)

	return &models.User{
		ID: userID,
	}, nil
}

func (c *Conn) InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.User, error) {
	logger := log.With(c.logger, "method", "CreateUser")

	if request.Firstname == "" || request.Password == "" {
		return nil, errors.New("name and password are required fields")
	}

	client := pb.NewUserServiceClient(c.conn)

	user, err := client.InsertUser(ctx, &pb.InsertUserRequest{
		Email:     request.Email,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Password:  request.Password,
		BirthDate: request.BirthDate,
		Address:   request.Address,
	})
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	userID, err := strconv.Atoi(user.ID)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	return &models.User{
		ID: userID,
	}, nil
}
