package friend

type EmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func ValidateEmailRequest(req *CreateFriend) error {
	return validate.Struct(req)
}
