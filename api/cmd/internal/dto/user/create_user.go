package user

type CreateUser struct {
	Email string `json:"requestor" validate:"required,email"`
}

func ValidateCreateUserRequest(req *CreateUser) error {
	return validate.Struct(req)
}
