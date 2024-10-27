package user

type CreateUser struct {
	Email string `json:"email" validate:"required,email"`
}

func ValidateCreateUserRequest(req *CreateUser) error {
	return validate.Struct(req)
}
