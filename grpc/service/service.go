package service

import (
	"context"
	"errors"

	"github.com/c95rt/bootcamp-user/grpc/entities"
	"github.com/c95rt/bootcamp-user/grpc/repository"
	"github.com/c95rt/bootcamp-user/helpers"
)

type service struct {
	repository repository.Repository
}

type Service interface {
	Login(ctx context.Context, request *entities.LoginRequest) (*entities.User, error)
	InsertUser(ctx context.Context, request *entities.InsertUserRequest) (*entities.User, error)
	/*UpdateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error*/
}

func NewService(rep repository.Repository) Service {
	return &service{
		repository: rep,
	}
}

func (s service) Login(ctx context.Context, request *entities.LoginRequest) (*entities.User, error) {
	// validate struct

	user, err := s.repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("email not found")
	}

	if !helpers.CheckPassword(user.Password, request.Password) {
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func (s service) InsertUser(ctx context.Context, request *entities.InsertUserRequest) (*entities.User, error) {
	// validate struct

	emailExists, err := s.repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if emailExists != nil {
		return nil, errors.New("email already exists")
	}

	request.Password, err = helpers.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.InsertUser(ctx, request)
	if err != nil {
		return nil, err
	}

	return user, nil

}
