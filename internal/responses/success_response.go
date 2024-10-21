package responses

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Success  bool        `json:"success"`
	Status   int         `json:"status"`
	MetaData interface{} `json:"metaData,omitempty"`
}

const (
	STATUS_OK      = http.StatusOK
	STATUS_CREATED = http.StatusCreated
)

var ReasonStatusCodeSuccess = map[int]string{
	STATUS_OK:      "Success",
	STATUS_CREATED: "Created",
}

func NewSuccessResponse(success bool, statusCode int, metaData interface{}) SuccessResponse {
	return SuccessResponse{
		Success:  success,
		Status:   statusCode,
		MetaData: metaData,
	}
}

func (sr *SuccessResponse) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sr.Status)
	json.NewEncoder(w).Encode(sr)
}

type OK struct {
	SuccessResponse
}

func NewOK(metaData interface{}) OK {
	return OK{
		SuccessResponse: NewSuccessResponse(true, STATUS_OK, metaData),
	}
}

type CREATED struct {
	SuccessResponse
}

func NewCREATED(metaData interface{}) CREATED {
	return CREATED{
		SuccessResponse: NewSuccessResponse(true, STATUS_CREATED, metaData),
	}
}
