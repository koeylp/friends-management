package friend

type CreateFriend struct {
	Friends []string `json:"friends" validate:"required,dive,email"`
}

func ValidateCreateFriendRequest(req *CreateFriend) error {
	return validate.Struct(req)
}
