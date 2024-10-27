package relationship

import (
	"context"
	"errors"

	"github.com/koeylp/friends-management/cmd/internal/model/dto/user"
	"github.com/stretchr/testify/mock"
)

// MockRelationshipRepository is a mock implementation of a relationship repository for testing purposes.
type MockRelationshipRepository struct {
	mock.Mock
}

// GetUpdatableEmailAddresses mocks the retrieval of email addresses that can be updated by a sender.
func (m *MockRelationshipRepository) GetUpdatableEmailAddresses(ctx context.Context, sender_id string) ([]string, error) {
	args := m.Called(ctx, sender_id)
	return args.Get(0).([]string), args.Error(1)
}

// BlockUpdates mocks the blocking of updates from a target user by a requestor.
func (m *MockRelationshipRepository) BlockUpdates(ctx context.Context, requestor_id string, target_id string) error {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Error(0)
}

// CheckBlockExists mocks the check for whether a block exists between two users.
func (m *MockRelationshipRepository) CheckBlockExists(ctx context.Context, requestor_id string, target_id string) (bool, error) {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Bool(0), args.Error(1)
}

// CheckSubscriptionExists mocks the check for whether a subscription exists between two users.
func (m *MockRelationshipRepository) CheckSubscriptionExists(ctx context.Context, requestor_id string, target_id string) (bool, error) {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Bool(0), args.Error(1)
}

// Subscribe mocks the subscription of a requestor to a target user.
func (m *MockRelationshipRepository) Subscribe(ctx context.Context, requestor_id string, target_id string) error {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Error(0)
}

// GetFriends mocks the retrieval of a list of friends for a given email address.
func (m *MockRelationshipRepository) GetFriends(ctx context.Context, email string) ([]string, error) {
	args := m.Called(ctx, email)
	return args.Get(0).([]string), args.Error(1)
}

// CreateFriend mocks the creation of a friendship between two users.
func (m *MockRelationshipRepository) CreateFriend(ctx context.Context, requestor_id, target_id string) error {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Error(0)
}

// CheckFriendshipExists mocks the existence of the frienship between 2 users
func (m *MockRelationshipRepository) CheckFriendshipExists(ctx context.Context, requestor_id, target_id string) (bool, error) {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Bool(0), args.Error(1)
}

// GetCommonFriends mocks the retrieval of common friends between a list of users.
func (m *MockRelationshipRepository) GetCommonFriends(ctx context.Context, users []*user.User) ([]string, error) {
	args := m.Called(ctx, users)
	return args.Get(0).([]string), args.Error(1)
}

// MockUserRepository is a mock implementation of a user repository for testing purposes.
type MockUserRepository struct {
	ShouldFail bool
	mock.Mock
}

// GetUserByEmail implements user.UserRepository.
func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

// CreateUser mocks the creation of a user.
// It returns an error if ShouldFail is set to true, otherwise it returns nil (success).
func (m *MockUserRepository) CreateUser(ctx context.Context, u *user.CreateUser) error {
	if m.ShouldFail {
		return errors.New("mock error: failed to create user")
	}
	return nil
}
