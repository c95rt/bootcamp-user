package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	"github.com/c95rt/bootcamp-user/http/endpoints"
	"github.com/c95rt/bootcamp-user/http/models"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPServer(endpoints endpoints.Endpoints, logger log.Logger) http.Handler {
	fmt.Println("Debug NewHTTPServer...")
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/login").Handler(
		httptransport.NewServer(
			endpoints.Login,
			decodeLoginRequest,
			encodeLoginResponse,
		),
	)

	r.Methods("POST").Path("/user").Handler(
		httptransport.NewServer(
			endpoints.InsertUser,
			decodeInsertUserRequest,
			encodeInsertUserResponse,
		),
	)

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func decodeLoginRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req models.LoginRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(models.User)
	return json.NewEncoder(w).Encode(res)
}

func decodeInsertUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req models.InsertUserRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeInsertUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(models.User)
	return json.NewEncoder(w).Encode(res)
}
