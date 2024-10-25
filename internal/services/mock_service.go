package services

import (
	"context"
	"errors"

	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	ShouldFail bool
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, u *user.CreateUser) error {
	if m.ShouldFail {
		return errors.New("mock error: failed to create user")
	}
	return nil
}

type MockRelationshipRepository struct {
	mock.Mock
}

func (m *MockRelationshipRepository) GetUpdatableEmailAddresses(ctx context.Context, sender_id string) ([]string, error) {
	args := m.Called(ctx, sender_id)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRelationshipRepository) BlockUpdates(ctx context.Context, requestor_id string, target_id string) error {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Error(0)
}

func (m *MockRelationshipRepository) CheckBlockExists(ctx context.Context, requestor_id string, target_id string) (bool, error) {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRelationshipRepository) CheckSubscriptionExists(ctx context.Context, requestor_id string, target_id string) (bool, error) {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRelationshipRepository) Subscribe(ctx context.Context, requestor_id string, target_id string) error {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Error(0)
}

func (m *MockRelationshipRepository) GetFriends(ctx context.Context, email string) ([]string, error) {
	args := m.Called(ctx, email)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRelationshipRepository) CreateFriend(ctx context.Context, requestor_id, target_id string) error {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Error(0)
}

func (m *MockRelationshipRepository) CheckFriendshipExists(ctx context.Context, requestor_id, target_id string) (bool, error) {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRelationshipRepository) GetCommonFriends(ctx context.Context, users []*user.User) ([]string, error) {
	args := m.Called(ctx, users)
	return args.Get(0).([]string), args.Error(1)
}
