package repository

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/c95rt/bootcamp-user/grpc/models"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, request *models.UpdateUserRequest) (*models.User, error)
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

func (db *DBConn) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	stmt, err := db.MariaDB.PrepareNamed(getUserByEmailQuery)
	if err != nil {
		return &models.User{}, err
	}
	args := map[string]interface{}{
		"email":  email,
		"active": DEFAULT_ACTIVE,
	}
	row := stmt.QueryRow(args)
	var user models.User
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
		return &models.User{}, err
	}
	return &user, nil
}

func (db *DBConn) InsertUser(ctx context.Context, request *models.InsertUserRequest) (*models.User, error) {
	tx, err := db.MariaDB.NewTx()
	if err != nil {
		return &models.User{}, errors.Wrap(err, "failed to start transaction")
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
		return &models.User{}, err
	}

	userAdditional, err := db.insertUserAdditionalTx(ctx, tx, id, request)
	if err != nil {
		return &models.User{}, err
	}

	return &models.User{
		ID:         id,
		Email:      request.Email,
		Firstname:  request.Firstname,
		Lastname:   request.Lastname,
		Password:   request.Password,
		Active:     DEFAULT_ACTIVE,
		Additional: userAdditional,
	}, nil
}

func (db *DBConn) insertUserTx(ctx context.Context, tx Tx, request *models.InsertUserRequest) (int, error) {
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

func (db *DBConn) insertUserAdditionalTx(ctx context.Context, tx Tx, id int, request *models.InsertUserRequest) (*models.UserAdditional, error) {
	additional := map[string]interface{}{
		"user_id":    id,
		"birth_date": request.BirthDate,
		"address":    request.Address,
	}

	_, err := db.MongoDB.Collection(USER_ADDITIONAL_COLLECTION).InsertOne(ctx, additional)
	if err != nil {
		return nil, err
	}

	return &models.UserAdditional{
		BirthDate: request.BirthDate,
		Address:   request.Address,
	}, nil
}

func (db *DBConn) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	stmt, err := db.MariaDB.PrepareNamed(getUserByIDQuery)
	if err != nil {
		return &models.User{}, err
	}

	args := map[string]interface{}{
		"user_id": userID,
		"active":  DEFAULT_ACTIVE,
	}

	row := stmt.QueryRow(args)

	var user models.User
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
		return &models.User{}, err
	}

	userAdditionalDocument := db.MongoDB.Collection(USER_ADDITIONAL_COLLECTION).FindOne(ctx, bson.D{{"user_id", userID}})
	if err := userAdditionalDocument.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			return &models.User{}, err
		}
	} else {
		userAdditional := models.UserAdditional{}
		userAdditionalDocument.Decode(&userAdditional)
		user.Additional = &userAdditional
	}

	return &user, nil
}

func (db *DBConn) UpdateUser(ctx context.Context, request *models.UpdateUserRequest) (*models.User, error) {
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
		return &models.User{}, err
	}

	return &models.User{
		ID:        request.ID,
		Email:     request.Email,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
	}, nil
}

func (db *DBConn) DeleteUser(ctx context.Context, userID int) error {
	tx, err := db.MariaDB.NewTx()
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	err = db.deleteUserTx(ctx, tx, userID)
	if err != nil {
		return err
	}

	err = db.deleteUserAdditional(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (db *DBConn) deleteUserTx(ctx context.Context, tx Tx, userID int) error {
	stmt, err := tx.PrepareNamed(deleteUserQuery)
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

func (db *DBConn) deleteUserAdditional(ctx context.Context, userID int) error {
	_, err := db.MongoDB.Collection(USER_ADDITIONAL_COLLECTION).DeleteMany(ctx, bson.D{{"user_id", userID}})
	if err != nil {
		return err
	}
	return nil
}
