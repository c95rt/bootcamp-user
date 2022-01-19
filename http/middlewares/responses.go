package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	STATUS_INTERNAL_SERVER_ERROR string = `internal server error`
	STATUS_NOT_FOUND_ERROR       string = `not found`
	STATUS_FORBIDDEN_ERROR       string = `forbidden`
	STATUS_UNAUTHORIZED_ERROR    string = `unauthorized`
	STATUS_BAD_REQUEST_ERROR     string = `bad request`
	STATUS_UNKNOWN_ERROR         string = `unknown error`
)

type ResponseWriter struct {
	writer http.ResponseWriter
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		writer: w,
	}
}

func (r *ResponseWriter) writePlainJSONResponse(statusCode int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		r.writer.WriteHeader(http.StatusInternalServerError)
		r.writer.Write([]byte(fmt.Sprintf("unexpected error: %v", err)))
	}

	r.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.writer.WriteHeader(statusCode)

	if code, err := r.writer.Write(b); err != nil {
		fmt.Sprintf("could not response - code: %d", code)
	}
}

func (er *ErrorResponse) newErrorResponse(statusCode int) {
	switch statusCode {
	case http.StatusInternalServerError:
		er.setErrorMessage(STATUS_INTERNAL_SERVER_ERROR)
		return
	case http.StatusNotFound:
		er.setErrorMessage(STATUS_NOT_FOUND_ERROR)
		return
	case http.StatusForbidden:
		er.setErrorMessage(STATUS_FORBIDDEN_ERROR)
		return
	case http.StatusUnauthorized:
		er.setErrorMessage(STATUS_UNAUTHORIZED_ERROR)
		return
	case http.StatusBadRequest:
		er.setErrorMessage(STATUS_BAD_REQUEST_ERROR)
		return
	default:
		er.setErrorMessage(STATUS_UNKNOWN_ERROR)
	}
}

func (er *ErrorResponse) setErrorMessage(message string) {
	er.Error = message
}

func (r *ResponseWriter) WriteJSON(statusCode int, data interface{}) {
	r.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !(statusCode >= 100 && statusCode <= 299) {
		errorResponse := ErrorResponse{}
		errorResponse.newErrorResponse(statusCode)
		r.writePlainJSONResponse(statusCode, errorResponse)
	}
	r.writePlainJSONResponse(statusCode, data)
}
