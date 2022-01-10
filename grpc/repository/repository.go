package repository

import (
	"github.com/jmoiron/sqlx"
)

type DBConn struct {
	MariaDB *DB
}

type Repository interface {
	UserRepository
}

func NewRepository(db *sqlx.DB) (Repository, error) {
	mariaDBConn, err := New(db)
	if err != nil {
		return nil, err
	}
	return &DBConn{
		MariaDB: mariaDBConn,
	}, nil
}
