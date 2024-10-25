package handlers

import (
	"context"

	"github.com/koeylp/friends-management/internal/dto/relationship/block"
	"github.com/koeylp/friends-management/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/internal/dto/relationship/subscription"
	"github.com/koeylp/friends-management/internal/dto/user"
)

type MockRelationshipService struct {
	CreateFriendFunc               func(ctx context.Context, req *friend.CreateFriend) error
	GetFriendListByEmailFunc       func(ctx context.Context, email string) ([]string, error)
	GetCommonListFunc              func(ctx context.Context, req *friend.CommonFriendListReq) ([]string, error)
	SubscribeFunc                  func(ctx context.Context, req *subscription.SubscribeRequest) error
	BlockUpdatesFunc               func(ctx context.Context, req *block.BlockRequest) error
	GetUpdatableEmailAddressesFunc func(ctx context.Context, req *subscription.RecipientRequest) ([]string, error)
}

type MockUserService struct {
	CreateUserFunc func(ctx context.Context, req *user.CreateUser) error
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	panic("unimplemented")
}

func (m *MockUserService) CreateUser(ctx context.Context, req *user.CreateUser) error {
	return m.CreateUserFunc(ctx, req)
}

func setupUserHandler(mockService *MockUserService) *UserHandler {
	return NewUserHandler(mockService)
}

func (m *MockRelationshipService) CreateFriend(ctx context.Context, req *friend.CreateFriend) error {
	return m.CreateFriendFunc(ctx, req)
}

func (m *MockRelationshipService) GetFriendListByEmail(ctx context.Context, email string) ([]string, error) {
	return m.GetFriendListByEmailFunc(ctx, email)
}

func (m *MockRelationshipService) GetCommonList(ctx context.Context, req *friend.CommonFriendListReq) ([]string, error) {
	return m.GetCommonListFunc(ctx, req)
}

func (m *MockRelationshipService) Subscribe(ctx context.Context, req *subscription.SubscribeRequest) error {
	return m.SubscribeFunc(ctx, req)
}

func (m *MockRelationshipService) BlockUpdates(ctx context.Context, req *block.BlockRequest) error {
	return m.BlockUpdatesFunc(ctx, req)
}

func (m *MockRelationshipService) GetUpdatableEmailAddresses(ctx context.Context, req *subscription.RecipientRequest) ([]string, error) {
	return m.GetUpdatableEmailAddressesFunc(ctx, req)
}

func setupRelationshipHandler(mockService *MockRelationshipService) *RelationshipHandler {
	return NewRelationshipHandler(mockService)
}
