package user

import (
	"context"
	"errors"

	"github.com/koeylp/friends-management/cmd/internal/model/dto/user"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of a user repository for testing purposes.
type MockUserRepository struct {
	ShouldFail bool
	mock.Mock
}

// CreateUser mocks the creation of a user.
// It returns an error if ShouldFail is set to true, otherwise it returns nil (success).
func (m *MockUserRepository) CreateUser(ctx context.Context, u *user.CreateUser) error {
	if m.ShouldFail {
		return errors.New("mock error: failed to create user")
	}
	return nil
}
