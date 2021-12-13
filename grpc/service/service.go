package service

import (
	"context"
	"errors"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/c95rt/bootcamp-user/grpc/repository"
	"github.com/c95rt/bootcamp-user/helpers"
)

type service struct {
	repository repository.Repository
	logger     log.Logger
}

type Service interface {
	Login(ctx context.Context, request *repository.LoginRequest) (*repository.User, error)
	InsertUser(ctx context.Context, request *repository.InsertUserRequest) (*repository.User, error)
	/*UpdateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error*/
}

func NewService(rep repository.Repository, logger log.Logger) Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

func (s service) Login(ctx context.Context, request *repository.LoginRequest) (*repository.User, error) {
	logger := log.With(s.logger, "method", "Login")

	// validate struct

	user, err := s.repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	if user == nil {
		level.Error(logger).Log("err", "email not found")
		return nil, errors.New("email not found")
	}

	if !helpers.CheckPassword(user.Password, request.Password) {
		level.Error(logger).Log("err", "Wrong password")
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func (s service) InsertUser(ctx context.Context, request *repository.InsertUserRequest) (*repository.User, error) {
	logger := log.With(s.logger, "method", "Api CreateUser")

	// validate struct

	emailExists, err := s.repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	if emailExists != nil {
		level.Error(logger).Log("err", "email already exists")
		return nil, errors.New("email already exists")
	}

	request.Password, err = helpers.HashPassword(request.Password)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	user, err := s.repository.InsertUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	level.Info(logger).Log("info", "User created")

	return user, nil

}
