package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/c95rt/bootcamp-user/http/models"
	"github.com/c95rt/bootcamp-user/http/repository"
)

type Service interface {
	Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error)
	InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.InsertUserResponse, error)
	GetUser(ctx context.Context, request *models.GetUserRequest) (*models.GetUserResponse, error)
	UpdateUser(ctx context.Context, request *models.UpdateUserRequest) (*models.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, request *models.DeleteUserRequest) (*models.DeleteUserResponse, error)
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

func (s service) Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error) {
	logger := log.With(s.logger, "method", "Login")

	response, err := s.repository.Login(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	level.Info(logger).Log("success", true)

	return response, nil
}

func (s service) InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.InsertUserResponse, error) {
	logger := log.With(s.logger, "method", "CreateUser")

	response, err := s.repository.InsertUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	level.Info(logger).Log("success", true)

	return response, nil
}

func (s service) UpdateUser(ctx context.Context, request *models.UpdateUserRequest) (*models.UpdateUserResponse, error) {
	logger := log.With(s.logger, "method", "UpdateUser")

	response, err := s.repository.UpdateUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	level.Info(logger).Log("success", true)

	return response, nil
}

func (s service) GetUser(ctx context.Context, request *models.GetUserRequest) (*models.GetUserResponse, error) {
	logger := log.With(s.logger, "method", "GetUser")

	response, err := s.repository.GetUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	level.Info(logger).Log("success", true)

	return response, nil
}

func (s service) DeleteUser(ctx context.Context, request *models.DeleteUserRequest) (*models.DeleteUserResponse, error) {
	logger := log.With(s.logger, "method", "DeleteUser")

	response, err := s.repository.DeleteUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	level.Info(logger).Log("success", true)

	return response, nil
}
