package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	jwtRequest "github.com/dgrijalva/jwt-go/request"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/c95rt/bootcamp-user/http/config"
	"github.com/c95rt/bootcamp-user/http/endpoints"
	httpErrors "github.com/c95rt/bootcamp-user/http/errors"
	"github.com/c95rt/bootcamp-user/http/middlewares"
	"github.com/c95rt/bootcamp-user/http/models"
)

func NewHTTPServer(appConfig *config.AppConfig, endpoints endpoints.Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(middlewares.CommonMiddleware)
	// r.Use(middlewares.JwtMiddleware([]byte(appConfig.Config.JWTSecret)))
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
	r.Methods("GET").Path("/user/{id:[0-9]+}").Handler(
		httptransport.NewServer(
			endpoints.GetUser,
			decodeGetUserRequest,
			encodeGetUserResponse,
		),
	)
	r.Methods("PUT").Path("/user/{id:[0-9]+}").Handler(
		httptransport.NewServer(
			endpoints.UpdateUser,
			decodeUpdateUserRequest,
			encodeUpdateUserResponse,
		),
	)
	r.Methods("DELETE").Path("/user/{id:[0-9]+}").Handler(
		httptransport.NewServer(
			endpoints.DeleteUser,
			decodeDeleteUserRequest,
			encodeDeleteUserResponse,
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

func jwtMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := jwtRequest.ParseFromRequest(r, jwtRequest.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			}); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.LoginRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, httpErrors.NewBadRequestError()
	}
	return req, nil
}

func encodeLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	rw := middlewares.NewResponseWriter(w)
	res, ok := response.(*models.LoginResponse)
	if !ok {
		rw.WriteJSON(http.StatusInternalServerError, nil)
		return nil
	}
	rw.WriteJSON(http.StatusOK, res)
	return nil
}

func decodeInsertUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.InsertUserRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, httpErrors.NewBadRequestError()
	}
	return req, nil
}

func encodeInsertUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(*models.InsertUserResponse)
	return json.NewEncoder(w).Encode(res)
}

func decodeUpdateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, httpErrors.NewBadRequestError()
	}
	vars := mux.Vars(r)
	req.ID, err = strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeUpdateUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(*models.UpdateUserResponse)
	return json.NewEncoder(w).Encode(res)
}

func decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.GetUserRequest
	var err error
	vars := mux.Vars(r)
	req.ID, err = strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeGetUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(*models.GetUserResponse)
	return json.NewEncoder(w).Encode(res)
}

func decodeDeleteUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.DeleteUserRequest
	var err error
	vars := mux.Vars(r)
	req.ID, err = strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeDeleteUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(*models.DeleteUserResponse)
	return json.NewEncoder(w).Encode(res)
}
