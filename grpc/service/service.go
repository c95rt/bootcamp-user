package service

import (
	"context"
	"errors"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/c95rt/bootcamp-user/grpc/config"
	grpcErrors "github.com/c95rt/bootcamp-user/grpc/errors"
	"github.com/c95rt/bootcamp-user/grpc/models"
	"github.com/c95rt/bootcamp-user/grpc/repository"
	"github.com/c95rt/bootcamp-user/helpers"
)

type service struct {
	repository repository.Repository
	logger     log.Logger
	appConfig  config.AppConfig
}

type Service interface {
	Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error)
	InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.InsertUserResponse, error)
	UpdateUser(ctx context.Context, request *models.UpdateUserRequest) (*models.UpdateUserResponse, error)
	GetUser(ctx context.Context, request *models.GetUserRequest) (*models.GetUserResponse, error)
	DeleteUser(ctx context.Context, request *models.DeleteUserRequest) (*models.DeleteUserResponse, error)
}

func NewService(rep repository.Repository, logger log.Logger) Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

func (s service) Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error) {
	logger := log.With(s.logger, "method", "GetUser")
	response := models.LoginResponse{}

	if request.Email == "" || request.Password == "" {
		level.Error(logger).Log("err", errors.New("failed validations"))
		return &response, grpcErrors.NewBadRequestError()
	}

	user, err := s.repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		level.Error(logger).Log("err", err)
		return &response, grpcErrors.NewInternalServerError()
	}

	if user == nil {
		level.Error(logger).Log("err", errors.New("email not found"))
		return &response, grpcErrors.NewNotFoundError()
	}

	if !helpers.CheckPassword(user.Password, request.Password) {
		level.Error(logger).Log("err", errors.New("bad password"))
		return &response, grpcErrors.NewBadRequestError()
	}

	response.Token, err = helpers.GenerateToken(user, s.appConfig.Config.JWTSecret)
	if err != nil {
		level.Error(logger).Log("err", err)
		return &response, grpcErrors.NewInternalServerError()
	}

	level.Info(logger).Log("success", true)

	return &response, nil
}

func (s service) InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.InsertUserResponse, error) {
	logger := log.With(s.logger, "method", "InsertUser")
	response := models.InsertUserResponse{}

	if request.Email == "" || request.Password == "" || request.Firstname == "" || request.Lastname == "" {
		level.Error(logger).Log("err", errors.New("failed validations"))
		return &response, grpcErrors.NewBadRequestError()
	}

	emailExists, err := s.repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		level.Error(logger).Log("err", err)
		return &response, grpcErrors.NewInternalServerError()
	}

	if emailExists != nil {
		level.Error(logger).Log("err", "email already exists")
		return &response, grpcErrors.NewBadRequestError()
	}

	request.Password, err = helpers.HashPassword(request.Password)
	if err != nil {
		level.Error(logger).Log("err", err)
		return &response, grpcErrors.NewInternalServerError()
	}

	user, err := s.repository.InsertUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return &response, grpcErrors.NewInternalServerError()
	}

	response.ID = user.ID
	response.Email = user.Email
	response.Firstname = user.Firstname
	response.Lastname = user.Lastname
	response.Password = user.Password

	if user.Additional != nil {
		response.BirthDate = user.Additional.BirthDate
		response.Address = user.Additional.Address
	}

	level.Info(logger).Log("success", true)

	return &response, nil
}

func (s service) UpdateUser(ctx context.Context, request *models.UpdateUserRequest) (*models.UpdateUserResponse, error) {
	logger := log.With(s.logger, "method", "UpdateUser")
	response := models.UpdateUserResponse{}

	var updateUser bool

	user, err := s.repository.GetUserByID(ctx, request.ID)
	if err != nil {
		level.Error(logger).Log("err", err)
		return &response, grpcErrors.NewInternalServerError()
	}

	if user == nil {
		level.Error(logger).Log("err", errors.New("user not found"))
		return &response, grpcErrors.NewNotFoundError()
	}

	if user.Email != request.Email {
		emailExists, err := s.repository.GetUserByEmail(ctx, request.Email)
		if err != nil {
			level.Error(logger).Log("err", err)
			return &response, grpcErrors.NewInternalServerError()
		}

		if emailExists != nil {
			level.Error(logger).Log("err", errors.New("email already exists"))
			return &response, grpcErrors.NewBadRequestError()
		}

		updateUser = true
	}

	request.Password, err = helpers.HashPassword(request.Password)
	if err != nil {
		level.Error(logger).Log("err", err)
		return &response, grpcErrors.NewInternalServerError()
	}
	if !helpers.CheckPassword(user.Password, request.Password) {
		updateUser = true
	}

	if user.Firstname != request.Firstname {
		updateUser = true
	}

	if user.Lastname != request.Lastname {
		updateUser = true
	}

	if updateUser {
		user, err = s.repository.UpdateUser(ctx, request)
		if err != nil {
			level.Error(logger).Log("err", err)
			return &response, grpcErrors.NewInternalServerError()
		}

		response.ID = user.ID
		response.Email = user.Email
		response.Firstname = user.Firstname
		response.Lastname = user.Lastname
		response.Password = user.Password

		if user.Additional != nil {
			response.BirthDate = user.Additional.BirthDate
			response.Address = user.Additional.Address
		}
	}

	level.Info(logger).Log("success", true)

	return &response, nil
}

func (s service) GetUser(ctx context.Context, request *models.GetUserRequest) (*models.GetUserResponse, error) {
	logger := log.With(s.logger, "method", "GetUser")

	user, err := s.repository.GetUserByID(ctx, request.ID)
	if err != nil {
		level.Error(logger).Log("err", err)
		return &models.GetUserResponse{}, grpcErrors.NewInternalServerError()
	}

	if user == nil {
		level.Error(logger).Log("err", errors.New("user not found"))
		return &models.GetUserResponse{}, grpcErrors.NewNotFoundError()
	}

	response := &models.GetUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Password:  user.Password,
	}

	if user.Additional != nil {
		response.BirthDate = user.Additional.BirthDate
		response.Address = user.Additional.Address
	}

	level.Info(logger).Log("success", true)

	return response, nil
}

func (s service) DeleteUser(ctx context.Context, request *models.DeleteUserRequest) (*models.DeleteUserResponse, error) {
	logger := log.With(s.logger, "method", "DeleteUser")

	user, err := s.repository.GetUserByID(ctx, request.ID)
	if err != nil {
		level.Error(logger).Log("err", err)
		return &models.DeleteUserResponse{}, grpcErrors.NewInternalServerError()
	}

	if user == nil {
		level.Error(logger).Log("err", errors.New("user not found"))
		return &models.DeleteUserResponse{}, grpcErrors.NewNotFoundError()
	}

	if err := s.repository.DeleteUser(ctx, user.ID); err != nil {
		level.Error(logger).Log("err", err)
		return &models.DeleteUserResponse{}, grpcErrors.NewInternalServerError()
	}

	level.Info(logger).Log("success", true)

	return &models.DeleteUserResponse{
		ID: user.ID,
	}, nil
}
