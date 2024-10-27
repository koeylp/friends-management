package subscription

type SubscribeRequest struct {
	Requestor string `json:"requestor" validate:"required,email"`
	Target    string `json:"target" validate:"required,email"`
}

func ValidateSubscribeRequest(req *SubscribeRequest) error {
	return validate.Struct(req)
}
