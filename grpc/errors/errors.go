package errors

import (
	"errors"
)

const (
	INTERNAL_SERVER_ERROR string = `internal_server_error`
	NOT_FOUND_ERROR       string = `not_found_error`
	BAD_REQUEST_ERROR     string = `bad_request_error`
)

func NewInternalServerError() error {
	return errors.New(INTERNAL_SERVER_ERROR)
}

func NewNotFoundError() error {
	return errors.New(NOT_FOUND_ERROR)
}

func NewBadRequestError() error {
	return errors.New(BAD_REQUEST_ERROR)
}
