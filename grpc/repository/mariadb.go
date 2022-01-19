package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type db interface {
	NewTx() (Tx, error)
}

type conn interface {
	Rebind(string) string
	NamedExec(string, interface{}) (sql.Result, error)
	Select(interface{}, string, ...interface{}) error
	QueryRow(string, ...interface{}) *sql.Row
	PrepareNamed(string) (*sqlx.NamedStmt, error)
	Get(interface{}, string, ...interface{}) error
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Preparex(query string) (*sqlx.Stmt, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

type Tx interface {
	conn

	Commit() error
	Rollback() error
}

type transactorImpl struct {
	*sqlx.DB
}

func (t *transactorImpl) NewTx() (Tx, error) {
	return t.Beginx()
}

type DB struct {
	conn
	db
}

func New(db *sqlx.DB) (*DB, error) {
	var (
		dbWrapper *DB
		err       error
	)

	tries := DEFAULT_MAX_RETRIES
	for tries >= 0 {
		time.Sleep(1 * time.Second)

		log.WithFields(log.Fields{
			"retries_left": tries,
		}).Warnf("%s: trying to connect to create connection", db.DriverName())

		dbWrapper, err = tryOpenConnection(db)
		if err != nil {
			if tries == 0 {
				return nil, err
			}

			tries = tries - 1
			continue
		}

		break
	}

	if err := createTableIfNotExists(db, createUserTableQuery); err != nil {
		return nil, err
	}

	return dbWrapper, nil
}

func tryOpenConnection(db *sqlx.DB) (*DB, error) {
	err := db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping db")
	}

	return &DB{
		db,
		&transactorImpl{db},
	}, nil
}

func createTableIfNotExists(db *sqlx.DB, query string) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
