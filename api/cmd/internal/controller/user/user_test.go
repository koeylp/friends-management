package user

import (
	"context"
	"errors"
	"testing"

	"github.com/koeylp/friends-management/cmd/internal/model/dto/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	userObj, ok := args.Get(0).(*user.User)
	if !ok {
		return nil, args.Error(1)
	}
	return userObj, args.Error(1)
}

// TestCreateUser_Success tests the successful creation of a user.
func TestCreateUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{ShouldFail: false}
	userController := NewUserController(mockRepo)

	sampleUser := &user.CreateUser{
		Email: "test@example.com",
	}

	err := userController.CreateUser(context.Background(), sampleUser)
	assert.NoError(t, err, "expected no error, got %v", err)
}

// TestCreateUser_Failure tests the failure scenario when creating a user.
func TestCreateUser_Failure(t *testing.T) {
	mockRepo := &MockUserRepository{ShouldFail: true}
	userController := NewUserController(mockRepo)

	sampleUser := &user.CreateUser{
		Email: "fail@example.com",
	}

	err := userController.CreateUser(context.Background(), sampleUser)

	assert.Error(t, err, "expected an error, got nil")
}

// TestGetUserByEmail tests retrieving a user by email successfully.
func TestGetUserByEmail(t *testing.T) {
	mockRepo := &MockUserRepository{ShouldFail: false}
	userController := NewUserController(mockRepo)

	expectedUser := &user.User{
		ID:    "123",
		Email: "test@example.com",
	}

	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(expectedUser, nil)

	result, err := userController.GetUserByEmail(context.Background(), "test@example.com")

	assert.NoError(t, err, "expected no error, got %v", err)
	assert.Equal(t, expectedUser, result, "expected and actual user mismatch")
}

// TestGetUserByEmail_NotFound tests retrieving a user that does not exist.
func TestGetUserByEmail_NotFound(t *testing.T) {
	mockRepo := &MockUserRepository{ShouldFail: false}
	userController := NewUserController(mockRepo)

	mockRepo.On("GetUserByEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New("user not found"))

	result, err := userController.GetUserByEmail(context.Background(), "notfound@example.com")

	assert.Error(t, err, "expected an error, got nil")
	assert.Nil(t, result, "expected nil user, got non-nil result")
}
