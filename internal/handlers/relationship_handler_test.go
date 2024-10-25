package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/koeylp/friends-management/internal/dto/relationship/block"
	"github.com/koeylp/friends-management/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/internal/dto/relationship/subscription"
	"github.com/stretchr/testify/assert"
)

// Test for creating a friendship relationship between two email addresses.
func TestCreateFriendHandler(t *testing.T) {
	mockService := &MockRelationshipService{
		CreateFriendFunc: func(ctx context.Context, req *friend.CreateFriend) error {
			return nil
		},
	}
	handler := setupRelationshipHandler(mockService)

	tests := []struct {
		name           string
		input          friend.CreateFriend
		expectedStatus int
	}{
		{
			name:           "Valid request",
			input:          friend.CreateFriend{Friends: []string{"user1@example.com", "user2@example.com"}},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid request - empty friends",
			input:          friend.CreateFriend{Friends: []string{}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid request - same friends",
			input:          friend.CreateFriend{Friends: []string{"user@example.com", "user@example.com"}},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/friends", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CreateFriendHandler(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

// Test for retrieving a list of friends for a specific email address.
func TestGetFriendListByEmailHandler(t *testing.T) {
	mockService := &MockRelationshipService{
		GetFriendListByEmailFunc: func(ctx context.Context, email string) ([]string, error) {
			return []string{"friend1@example.com", "friend2@example.com"}, nil
		},
	}
	handler := setupRelationshipHandler(mockService)

	tests := []struct {
		name            string
		input           friend.EmailRequest
		expectedStatus  int
		expectedFriends []string
	}{
		{
			name:            "Valid request",
			input:           friend.EmailRequest{Email: "user@example.com"},
			expectedStatus:  http.StatusOK,
			expectedFriends: []string{"friend1@example.com", "friend2@example.com"},
		},
		{
			name:           "Invalid request - empty email",
			input:          friend.EmailRequest{Email: ""},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/friends/list", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.GetFriendListByEmailHandler(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				json.NewDecoder(res.Body).Decode(&response)
				sort.Strings(tt.expectedFriends)
				friends := response["friends"].([]interface{})
				actualFriends := make([]string, len(friends))
				for i, friend := range friends {
					actualFriends[i] = friend.(string)
				}
				sort.Strings(actualFriends)
				assert.Equal(t, tt.expectedFriends, actualFriends)
			}
		})
	}
}

// Test for retrieving a common list of friends between two users.
func TestGetCommonListHandler(t *testing.T) {
	mockService := &MockRelationshipService{
		GetCommonListFunc: func(ctx context.Context, req *friend.CommonFriendListReq) ([]string, error) {
			return []string{"commonfriend@example.com"}, nil
		},
	}
	handler := setupRelationshipHandler(mockService)

	tests := []struct {
		name            string
		input           friend.CommonFriendListReq
		expectedStatus  int
		expectedFriends []string
	}{
		{
			name:            "Valid request",
			input:           friend.CommonFriendListReq{Friends: []string{"user1@example.com", "user2@example.com"}},
			expectedStatus:  http.StatusOK,
			expectedFriends: []string{"commonfriend@example.com"},
		},
		{
			name:           "Invalid request - same friends",
			input:          friend.CommonFriendListReq{Friends: []string{"user@example.com", "user@example.com"}},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/friends/common", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.GetCommonListHandler(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				json.NewDecoder(res.Body).Decode(&response)

				actualFriends := make([]string, len(response["friends"].([]interface{})))
				for i, friend := range response["friends"].([]interface{}) {
					actualFriends[i] = friend.(string)
				}

				sort.Strings(tt.expectedFriends)
				sort.Strings(actualFriends)

				assert.Equal(t, tt.expectedFriends, actualFriends)
			}
		})
	}
}

// Test for handling subscription requests.
func TestSubscribeHandler(t *testing.T) {
	mockService := &MockRelationshipService{
		SubscribeFunc: func(ctx context.Context, req *subscription.SubscribeRequest) error {
			return nil
		},
	}
	handler := setupRelationshipHandler(mockService)

	tests := []struct {
		name           string
		input          subscription.SubscribeRequest
		expectedStatus int
	}{
		{
			name:           "Valid subscription request",
			input:          subscription.SubscribeRequest{Requestor: "user@example.com", Target: "friend@example.com"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid subscription - same requestor and target",
			input:          subscription.SubscribeRequest{Requestor: "user@example.com", Target: "user@example.com"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON payload",
			input:          subscription.SubscribeRequest{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.SubscribeHandler(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

// Test for handling block updates requests.
func TestBlockUpdatesHandler(t *testing.T) {
	mockService := &MockRelationshipService{
		BlockUpdatesFunc: func(ctx context.Context, req *block.BlockRequest) error {
			return nil
		},
	}
	handler := setupRelationshipHandler(mockService)

	tests := []struct {
		name           string
		input          block.BlockRequest
		expectedStatus int
	}{
		{
			name:           "Valid block request",
			input:          block.BlockRequest{Requestor: "user@example.com", Target: "blockfriend@example.com"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid block - same requestor and target",
			input:          block.BlockRequest{Requestor: "user@example.com", Target: "user@example.com"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON payload",
			input:          block.BlockRequest{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/block", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.BlockUpdatesHandler(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

// Test for retrieving updatable email addresses based on sender's updates.
func TestGetUpdatableEmailAddressesHandler(t *testing.T) {
	mockService := &MockRelationshipService{
		GetUpdatableEmailAddressesFunc: func(ctx context.Context, req *subscription.RecipientRequest) ([]string, error) {
			return []string{"recipient1@example.com", "recipient2@example.com"}, nil
		},
	}
	handler := setupRelationshipHandler(mockService)

	tests := []struct {
		name           string
		input          subscription.RecipientRequest
		expectedStatus int
		expectedEmails []string
	}{
		{
			name:           "Valid request for updatable emails",
			input:          subscription.RecipientRequest{Sender: "user@example.com", Text: "Hello Wolrd! kate@example.com"},
			expectedStatus: http.StatusOK,
			expectedEmails: []string{"recipient1@example.com", "recipient2@example.com"},
		},
		{
			name:           "Invalid request - missing sender and text",
			input:          subscription.RecipientRequest{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/recipients", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.GetUpdatableEmailAddressesHandler(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.NewDecoder(res.Body).Decode(&response)
				assert.NoError(t, err)

				actualEmails := response["recipients"].([]interface{})
				var actualEmailsStr []string
				for _, email := range actualEmails {
					actualEmailsStr = append(actualEmailsStr, email.(string))
				}

				assert.Equal(t, tt.expectedEmails, actualEmailsStr)
			}
		})
	}
}
