package repository

import (
	"github.com/jmoiron/sqlx"

	"github.com/c95rt/bootcamp-user/grpc/repository/mariadb"
	"github.com/go-kit/log"
)

type DBConn struct {
	MariaDB *mariadb.DB
	logger  log.Logger
}

type Repository interface {
	UserRepository
}

func NewRepository(db *sqlx.DB, logger log.Logger) (Repository, error) {
	mariaDBConn, err := mariadb.New(db)
	if err != nil {
		return nil, err
	}
	return &DBConn{
		MariaDB: mariaDBConn,
		logger:  log.With(logger, "error", "db"),
	}, nil
}
