package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserHandler(t *testing.T) {
	mockService := &MockUserService{
		CreateUserFunc: func(ctx context.Context, req *user.CreateUser) error {
			return nil
		},
	}

	handler := setupUserHandler(mockService)

	tests := []struct {
		name           string
		input          user.CreateUser
		expectedStatus int
	}{
		{
			name:           "Valid user creation request",
			input:          user.CreateUser{Email: "user@example.com"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid request - JSON decoding error",
			input:          user.CreateUser{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "User creation failure",
			input:          user.CreateUser{Email: "user@example.com"},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			if tt.name == "User creation failure" {
				mockService.CreateUserFunc = func(ctx context.Context, req *user.CreateUser) error {
					return assert.AnError
				}
			} else {
				mockService.CreateUserFunc = func(ctx context.Context, req *user.CreateUser) error {
					return nil
				}
			}

			handler.CreateUserHandler(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				err := json.NewDecoder(res.Body).Decode(&response)
				assert.NoError(t, err)
			}
		})
	}
}
