package responses

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type SuccessResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"metaData,omitempty"`
	Status  int                    `json:"-"`
}

const (
	STATUS_OK      = http.StatusOK
	STATUS_CREATED = http.StatusCreated
)

func NewSuccessResponse(success bool, statusCode int, data map[string]interface{}) SuccessResponse {
	return SuccessResponse{
		Success: success,
		Status:  statusCode,
		Data:    data,
	}
}

func (sr *SuccessResponse) Send(w http.ResponseWriter) {
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(sr.Status)
	// json.NewEncoder(w).Encode(sr)

	response := make(map[string]interface{})

	response["count"] = len(sr.Data)
	for key, value := range sr.Data {
		if reflect.ValueOf(value).IsNil() {
			response[key] = []interface{}{}
			break
		}
		response[key] = value
	}

	response["success"] = sr.Success

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sr.Status)
	json.NewEncoder(w).Encode(response)
}

type OK struct {
	SuccessResponse
}

func NewOK(data map[string]interface{}) OK {
	return OK{
		SuccessResponse: NewSuccessResponse(true, STATUS_OK, data),
	}
}

type CREATED struct {
	SuccessResponse
}

func NewCREATED(data map[string]interface{}) CREATED {
	return CREATED{
		SuccessResponse: NewSuccessResponse(true, STATUS_CREATED, data),
	}
}
