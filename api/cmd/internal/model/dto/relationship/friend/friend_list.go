package friend

type FriendList struct {
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

func ValidateFriendListRequest(req *CreateFriend) error {
	return validate.Struct(req)
}
