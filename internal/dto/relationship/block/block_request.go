package block

import "github.com/go-playground/validator"

type BlockRequest struct {
	Requestor string `json:"requestor" validate:"required,email"`
	Target    string `json:"target" validate:"required,email"`
}

var validate = validator.New()

func ValidateBlockRequest(req *BlockRequest) error {
	return validate.Struct(req)
}
