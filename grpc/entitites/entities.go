package entity

type User struct {
	ID        int
	Email     string
	Firstname string
	Lastname  string
	Password  string
	Token     string
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

type UpdateUserRequest struct {
	ID        int
	Email     string
	Firstname string
	Lastname  string
	Password  string
}

type UserAdditional struct {
	BirthDate string
	Address   string
}

type LoginRequest struct {
	Email    string
	Password string
}
