package block

type BlockRequest struct {
	Requestor string `json:"requestor" validate:"required,email"`
	Target    string `json:"target" validate:"required,email"`
}

func ValidateBlockRequest(req *BlockRequest) error {
	return validate.Struct(req)
}
