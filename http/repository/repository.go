package repository

import (
	"github.com/go-kit/log"
	"google.golang.org/grpc"
)

type Conn struct {
	conn   *grpc.ClientConn
	logger log.Logger
}

type Repository interface {
	UserRepository
}

func NewRepository(conn *grpc.ClientConn) UserRepository {
	return &Conn{
		conn: conn,
	}
}
