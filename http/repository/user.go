package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/c95rt/bootcamp-user/grpc/pb"
	"github.com/c95rt/bootcamp-user/http/models"
)

type UserRepository interface {
	Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error)
	InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.InsertUserResponse, error)
	UpdateUser(ctx context.Context, request *models.UpdateUserRequest) (*models.UpdateUserResponse, error)
	GetUser(ctx context.Context, request *models.GetUserRequest) (*models.GetUserResponse, error)
	DeleteUser(ctx context.Context, request *models.DeleteUserRequest) (*models.DeleteUserResponse, error)
}

func (c *Conn) Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error) {
	client := pb.NewUserServiceClient(c.conn)

	response, err := client.Login(ctx, &pb.LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: response.Token,
	}, nil
}

func (c *Conn) InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.InsertUserResponse, error) {
	if request.Firstname == "" || request.Password == "" {
		return nil, errors.New("name and password are required fields")
	}

	client := pb.NewUserServiceClient(c.conn)

	response, err := client.InsertUser(ctx, &pb.InsertUserRequest{
		Email:     request.Email,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Password:  request.Password,
		BirthDate: request.BirthDate,
		Address:   request.Address,
	})
	if err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(response.Id)
	if err != nil {
		return nil, err
	}

	return &models.InsertUserResponse{
		ID:        userID,
		Email:     response.Email,
		Firstname: response.Firstname,
		Lastname:  response.Lastname,
		Password:  response.Password,
		BirthDate: response.BirthDate,
		Address:   response.Address,
	}, nil
}

func (c *Conn) UpdateUser(ctx context.Context, request *models.UpdateUserRequest) (*models.UpdateUserResponse, error) {
	client := pb.NewUserServiceClient(c.conn)

	response, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:        strconv.Itoa(request.ID),
		Email:     request.Email,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Password:  request.Password,
		BirthDate: request.BirthDate,
		Address:   request.Address,
	})
	if err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(response.Id)
	if err != nil {
		return nil, err
	}

	return &models.UpdateUserResponse{
		ID:        userID,
		Email:     response.Email,
		Firstname: response.Firstname,
		Lastname:  response.Lastname,
		Password:  response.Password,
		BirthDate: response.BirthDate,
		Address:   response.Address,
	}, nil
}

func (c *Conn) GetUser(ctx context.Context, request *models.GetUserRequest) (*models.GetUserResponse, error) {
	client := pb.NewUserServiceClient(c.conn)

	response, err := client.GetUser(ctx, &pb.GetUserRequest{
		Id: strconv.Itoa(request.ID),
	})
	if err != nil {
		return &models.GetUserResponse{}, err
	}

	userID, err := strconv.Atoi(response.Id)
	if err != nil {
		return &models.GetUserResponse{}, err
	}

	return &models.GetUserResponse{
		ID:        userID,
		Email:     response.Email,
		Firstname: response.Firstname,
		Lastname:  response.Lastname,
		Password:  response.Password,
		BirthDate: response.BirthDate,
		Address:   response.Address,
	}, nil
}

func (c *Conn) DeleteUser(ctx context.Context, request *models.DeleteUserRequest) (*models.DeleteUserResponse, error) {
	client := pb.NewUserServiceClient(c.conn)

	response, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: strconv.Itoa(request.ID),
	})
	if err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(response.Id)
	if err != nil {
		return nil, err
	}

	return &models.DeleteUserResponse{
		ID: userID,
	}, nil
}
