package friend

import (
	"github.com/go-playground/validator"
)

type CommonFriendListReq struct {
	Friends []string `json:"friends" validate:"required,dive,email"`
}

var validate = validator.New()

func ValidateCommonFriendListRequest(req *CommonFriendListReq) error {
	return validate.Struct(req)
}
