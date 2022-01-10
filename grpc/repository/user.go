package repository

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/c95rt/bootcamp-user/grpc/entities"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	InsertUser(ctx context.Context, request *entities.InsertUserRequest) (*entities.User, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
	UpdateUser(ctx context.Context, request *entities.UpdateUserRequest) (*entities.User, error)
	DeleteUser(ctx context.Context, id int) error
}

const (
	getUserByEmailQuery = `
		SELECT
			user.id,
			user.email,
			user.firstname,
			user.lastname,
			user.password,
			user.active
		FROM
			user
		WHERE
			user.email = :email AND
			user.active = :active
	`
	insertUserQuery = `
		INSERT
			user
		SET
			email = :email,
			firstname = :firstname,
			lastname = :lastname,
			password = :password
	`
	getUserByIDQuery = `
		SELECT
			user.id,
			user.email,
			user.firstname,
			user.lastname,
			user.active
		FROM
			user
		WHERE
			user.id = :user_id AND
			user.active = :active
	`
	updateUserQuery = `
		UPDATE
			user
		SET
			user.email = :email,
			user.firstname = :firstname,
			user.lastname = :lastname,
			user.password = :password
		WHERE
			user.id = :user_id
	`
	deleteUserQuery = `
		UPDATE
			user
		SET
			user.active = :inactive
		WHERE
			user.id = :user_id
	`
)

func (db *DBConn) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	// logger := log.With(repo.logger, "method", "Authenticate")

	stmt, err := db.MariaDB.PrepareNamed(getUserByEmailQuery)
	if err != nil {
		return nil, err
	}
	args := map[string]interface{}{
		"email":  email,
		"active": DEFAULT_ACTIVE,
	}
	row := stmt.QueryRow(args)
	var user entities.User
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Firstname,
		&user.Lastname,
		&user.Password,
		&user.Active,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (db *DBConn) InsertUser(ctx context.Context, request *entities.InsertUserRequest) (*entities.User, error) {
	tx, err := db.MariaDB.NewTx()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	id, err := db.insertUserTx(ctx, tx, request)
	if err != nil {
		return nil, err
	}

	userAdditional, err := db.insertUserAdditionalTx(ctx, tx, request)
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:         id,
		Email:      request.Email,
		Firstname:  request.Firstname,
		Lastname:   request.Lastname,
		Password:   request.Password,
		Active:     DEFAULT_ACTIVE,
		Additional: userAdditional,
	}, nil
}

func (db *DBConn) insertUserTx(ctx context.Context, tx Tx, request *entities.InsertUserRequest) (int, error) {
	stmt, err := db.MariaDB.PrepareNamed(insertUserQuery)
	if err != nil {
		return 0, err
	}

	args := map[string]interface{}{
		"email":     request.Email,
		"firstname": request.Firstname,
		"lastname":  request.Lastname,
		"password":  request.Password,
	}

	result, err := stmt.Exec(args)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (db *DBConn) insertUserAdditionalTx(ctx context.Context, tx Tx, request *entities.InsertUserRequest) (*entities.UserAdditional, error) {
	return &entities.UserAdditional{}, nil
}

func (db *DBConn) GetUserByID(ctx context.Context, userID int) (*entities.User, error) {
	stmt, err := db.MariaDB.PrepareNamed(getUserByIDQuery)
	if err != nil {
		return nil, err
	}

	args := map[string]interface{}{
		"user_id": userID,
		"active":  DEFAULT_ACTIVE,
	}

	row := stmt.QueryRow(args)

	var user entities.User
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Firstname,
		&user.Lastname,
		&user.Active,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (db *DBConn) UpdateUser(ctx context.Context, request *entities.UpdateUserRequest) (*entities.User, error) {
	stmt, err := db.MariaDB.PrepareNamed(updateUserQuery)
	if err != nil {
		return nil, err
	}

	args := map[string]interface{}{
		"user_id":   request.ID,
		"email":     request.Email,
		"firstname": request.Firstname,
		"lastname":  request.Lastname,
		"password":  request.Password,
	}

	_, err = stmt.Exec(args)
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:        request.ID,
		Email:     request.Email,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
	}, nil
}

func (db *DBConn) DeleteUser(ctx context.Context, userID int) error {
	stmt, err := db.MariaDB.PrepareNamed(deleteUserQuery)
	if err != nil {
		return err
	}

	args := map[string]interface{}{
		"user_id":  userID,
		"inactive": DEFAULT_INACTIVE,
	}

	_, err = stmt.Exec(args)
	if err != nil {
		return err
	}

	return nil
}
