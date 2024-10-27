package subscription

type RecipientRequest struct {
	Sender string `json:"sender" validate:"required,email"`
	Text   string `json:"text" validate:"required"`
}

func ValidateRecipientRequest(req *RecipientRequest) error {
	return validate.Struct(req)
}
