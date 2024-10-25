package friend

type EmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func ValidateEmailRequest(req *EmailRequest) error {
	return validate.Struct(req)
}
