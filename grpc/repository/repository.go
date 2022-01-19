package repository

import (
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBConn struct {
	MariaDB *DB
	MongoDB *mongo.Database
}

type Repository interface {
	UserRepository
}

func NewRepository(mariaDB *sqlx.DB, mongoDB *mongo.Database) (Repository, error) {
	mariaDBConn, err := New(mariaDB)
	if err != nil {
		return nil, err
	}
	return &DBConn{
		MariaDB: mariaDBConn,
		MongoDB: mongoDB,
	}, nil
}
