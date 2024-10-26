package handlers

import (
	"context"

	"github.com/koeylp/friends-management/cmd/internal/dto/relationship/block"
	"github.com/koeylp/friends-management/cmd/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/cmd/internal/dto/relationship/subscription"
	"github.com/koeylp/friends-management/cmd/internal/dto/user"
)

// MockRelationshipService is a mock implementation of a relationship service for testing purposes.
// It allows defining custom behaviors for its methods using function types.
type MockRelationshipService struct {
	CreateFriendFunc               func(ctx context.Context, req *friend.CreateFriend) error
	GetFriendListByEmailFunc       func(ctx context.Context, email string) ([]string, error)
	GetCommonListFunc              func(ctx context.Context, req *friend.CommonFriendListReq) ([]string, error)
	SubscribeFunc                  func(ctx context.Context, req *subscription.SubscribeRequest) error
	BlockUpdatesFunc               func(ctx context.Context, req *block.BlockRequest) error
	GetUpdatableEmailAddressesFunc func(ctx context.Context, req *subscription.RecipientRequest) ([]string, error)
}

// MockUserService is a mock implementation of a user service for testing purposes.
// It allows defining custom behavior for user-related methods.
type MockUserService struct {
	CreateUserFunc func(ctx context.Context, req *user.CreateUser) error
}

// setupRelationshipHandler initializes a RelationshipHandler with the provided mock relationship service.
func setupRelationshipHandler(mockService *MockRelationshipService) *RelationshipHandler {
	return NewRelationshipHandler(mockService)
}

// CreateFriend calls the custom CreateFriendFunc defined in the MockRelationshipService.
func (m *MockRelationshipService) CreateFriend(ctx context.Context, req *friend.CreateFriend) error {
	return m.CreateFriendFunc(ctx, req)
}

// GetUserByEmail is a method that needs to be implemented in the MockUserService.
// Currently, it panics if called, indicating it's not intended for use in this mock.
func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	panic("unimplemented")
}

// CreateUser calls the custom CreateUserFunc defined in the MockUserService.
func (m *MockUserService) CreateUser(ctx context.Context, req *user.CreateUser) error {
	return m.CreateUserFunc(ctx, req)
}

// setupUserHandler initializes a UserHandler with the provided mock user service.
func setupUserHandler(mockService *MockUserService) *UserHandler {
	return NewUserHandler(mockService)
}

// GetFriendListByEmail calls the custom GetFriendListByEmailFunc defined in the MockRelationshipService.
func (m *MockRelationshipService) GetFriendListByEmail(ctx context.Context, email string) ([]string, error) {
	return m.GetFriendListByEmailFunc(ctx, email)
}

// GetCommonList calls the custom GetCommonListFunc defined in the MockRelationshipService.
func (m *MockRelationshipService) GetCommonList(ctx context.Context, req *friend.CommonFriendListReq) ([]string, error) {
	return m.GetCommonListFunc(ctx, req)
}

// Subscribe calls the custom SubscribeFunc defined in the MockRelationshipService.
func (m *MockRelationshipService) Subscribe(ctx context.Context, req *subscription.SubscribeRequest) error {
	return m.SubscribeFunc(ctx, req)
}

// BlockUpdates calls the custom BlockUpdatesFunc defined in the MockRelationshipService.
func (m *MockRelationshipService) BlockUpdates(ctx context.Context, req *block.BlockRequest) error {
	return m.BlockUpdatesFunc(ctx, req)
}

// GetUpdatableEmailAddresses calls the custom GetUpdatableEmailAddressesFunc defined in the MockRelationshipService.
func (m *MockRelationshipService) GetUpdatableEmailAddresses(ctx context.Context, req *subscription.RecipientRequest) ([]string, error) {
	return m.GetUpdatableEmailAddressesFunc(ctx, req)
}
