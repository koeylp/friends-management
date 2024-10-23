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
	response := make(map[string]interface{})

	for key, value := range sr.Data {
		if reflect.ValueOf(value).IsNil() {
			response[key] = []interface{}{}
			response["count"] = 0
			break
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice, reflect.Array:
			slice := reflect.ValueOf(value)
			convertedSlice := make([]interface{}, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				convertedSlice[i] = slice.Index(i).Interface()
			}

			response[key] = convertedSlice
			response["count"] = len(convertedSlice)
		default:
			response[key] = value
		}
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
