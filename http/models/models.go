package models

type InsertUserRequest struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	BirthDate string `json:"birth_date"`
	Address   string `json:"address"`
}

type InsertUserResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	BirthDate string `json:"birth_date"`
	Address   string `json:"address"`
}

type UpdateUserRequest struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	BirthDate string `json:"birth_date"`
	Address   string `json:"address"`
}

type UpdateUserResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	BirthDate string `json:"birth_date"`
	Address   string `json:"address"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type DeleteUserRequest struct {
	ID int `json:"id"`
}

type DeleteUserResponse struct {
	ID int `json:"id"`
}

type GetUserRequest struct {
	ID int `json:"id"`
}

type GetUserResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	BirthDate string `json:"birth_date"`
	Address   string `json:"address"`
}
