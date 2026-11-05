package dto

type CreateUser struct {
	Email     string
	Password  string
	FirstName string
	LastName  *string
}

type UpdateUser struct {
	FirstName *string
	LastName  *string
}

type ChangePassword struct {
	Password string
}
