package models

type User struct {
	ID        int
	Email     string
	Firstname string
	Lastname  string
	Password  string
	Active    bool

	Additional *UserAdditional
}

type InsertUserRequest struct {
	Email     string
	Firstname string
	Lastname  string
	Password  string
	BirthDate string
	Address   string
}

type InsertUserResponse struct {
	ID        int
	Email     string
	Firstname string
	Lastname  string
	Password  string
	BirthDate string
	Address   string
}

type UpdateUserRequest struct {
	ID        int
	Email     string
	Firstname string
	Lastname  string
	Password  string
	BirthDate string
	Address   string
}

type UpdateUserResponse struct {
	ID        int
	Email     string
	Firstname string
	Lastname  string
	Password  string
	BirthDate string
	Address   string
}

type UserAdditional struct {
	BirthDate string `bson:"birth_date,omitempty"`
	Address   string `bson:"address,omitempty"`
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string
}

type DeleteUserRequest struct {
	ID int
}

type DeleteUserResponse struct {
	ID int
}

type GetUserRequest struct {
	ID int
}

type GetUserResponse struct {
	ID        int
	Email     string
	Firstname string
	Lastname  string
	Password  string
	BirthDate string
	Address   string
}
