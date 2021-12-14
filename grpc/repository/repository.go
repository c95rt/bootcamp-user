package repository

import (
	"github.com/jmoiron/sqlx"

	"github.com/c95rt/bootcamp-user/grpc/repository/mariadb"
)

type DBConn struct {
	MariaDB *mariadb.DB
}

type Repository interface {
	UserRepository
}

func NewRepository(db *sqlx.DB) (Repository, error) {
	mariaDBConn, err := mariadb.New(db)
	if err != nil {
		return nil, err
	}
	return &DBConn{
		MariaDB: mariaDBConn,
	}, nil
}
