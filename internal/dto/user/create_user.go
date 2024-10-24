package user

import "github.com/go-playground/validator"

type CreateUser struct {
	Email string `json:"requestor" validate:"required,email"`
}

var validate = validator.New()

func ValidateCreateUserRequest(req *CreateUser) error {
	return validate.Struct(req)
}
