package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/c95rt/bootcamp-user/http/models"
	"github.com/c95rt/bootcamp-user/http/repository"
)

type Service interface {
	Login(ctx context.Context, request *models.LoginRequest) (*models.User, error)
	InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.User, error)
}

type service struct {
	repository repository.UserRepository
	logger     log.Logger
}

func NewService(rep repository.UserRepository, logger log.Logger) Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

func (s service) Login(ctx context.Context, request *models.LoginRequest) (*models.User, error) {
	logger := log.With(s.logger, "method", "Authenticate")

	// TODO: validate struct

	user, err := s.repository.Login(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	return user, nil
}

func (s service) InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.User, error) {
	logger := log.With(s.logger, "method", "CreateUser")

	// TODO: validate struct

	user, err := s.repository.InsertUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	logger.Log("User created", user.ID)

	return user, nil
}
