package services

import (
	"context"
	"errors"
	"testing"

	"github.com/koeylp/friends-management/internal/dto/user"
)

type MockUserRepository struct {
	ShouldFail bool
}

func (m *MockUserRepository) CreateUser(ctx context.Context, u *user.CreateUser) error {
	if m.ShouldFail {
		return errors.New("mock error: failed to create user")
	}
	return nil
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	if m.ShouldFail {
		return nil, errors.New("mock error: failed to get user")
	}
	return nil, nil
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{ShouldFail: false}
	userService := NewUserService(mockRepo)

	sampleUser := &user.CreateUser{
		Email: "test@example.com",
	}

	err := userService.CreateUser(context.Background(), sampleUser)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreateUser_Failure(t *testing.T) {
	mockRepo := &MockUserRepository{ShouldFail: true}
	userService := NewUserService(mockRepo)

	sampleUser := &user.CreateUser{
		Email: "fail@example.com",
	}

	err := userService.CreateUser(context.Background(), sampleUser)

	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}
