package responses

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Message  string      `json:"message"`
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

func NewSuccessResponse(message string, statusCode int, metaData interface{}) SuccessResponse {

	if message == "" {
		message = ReasonStatusCodeSuccess[statusCode]
	}

	return SuccessResponse{
		Message:  message,
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

func NewOK(message string, metaData interface{}) OK {
	return OK{
		SuccessResponse: NewSuccessResponse(message, STATUS_OK, metaData),
	}
}

type CREATED struct {
	SuccessResponse
}

func NewCREATED(message string, metaData interface{}) CREATED {
	return CREATED{
		SuccessResponse: NewSuccessResponse(message, STATUS_CREATED, metaData),
	}
}
