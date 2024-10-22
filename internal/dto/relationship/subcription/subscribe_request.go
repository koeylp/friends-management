package subcription

import "github.com/go-playground/validator"

type SubscribeRequest struct {
	Requestor string `json:"requestor" validate:"required,email"`
	Target    string `json:"target" validate:"required,email"`
}

var validate = validator.New()

func ValidateSubscribeRequest(req *SubscribeRequest) error {
	return validate.Struct(req)
}
