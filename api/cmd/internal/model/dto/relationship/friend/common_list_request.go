package friend

type CommonFriendListReq struct {
	Friends []string `json:"friends" validate:"required,dive,email"`
}

func ValidateCommonFriendListRequest(req *CommonFriendListReq) error {
	return validate.Struct(req)
}
