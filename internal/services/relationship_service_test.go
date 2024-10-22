package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/koeylp/friends-management/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/koeylp/friends-management/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRelationshipRepository struct {
	mock.Mock
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

func TestCreateFriend(t *testing.T) {
	ctx := context.Background()

	mockRelRepo := new(MockRelationshipRepository)
	mockUserRepo := new(services.MockUserRepository)

	service := services.NewRelationshipService(mockRelRepo, mockUserRepo)

	inputEmails := []string{"requestor@example.com", "target@example.com"}
	input := &friend.CreateFriend{
		Friends: inputEmails,
	}

	mockUsers := []*user.User{
		{ID: "1", Email: "requestor@example.com"},
		{ID: "2", Email: "target@example.com"},
	}

	mockUserRepo.On("GetUserByEmail", ctx, "requestor@example.com").Return(mockUsers[0], nil)
	mockUserRepo.On("GetUserByEmail", ctx, "target@example.com").Return(mockUsers[1], nil)

	mockRelRepo.On("CheckFriendshipExists", ctx, "1", "2").Return(true, nil)

	err := service.CreateFriend(ctx, input)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "friendship already exists between requestor@example.com and target@example.com")

	mockRelRepo.ExpectedCalls = nil
	mockRelRepo.On("CheckFriendshipExists", ctx, "1", "2").Return(false, nil)
	mockRelRepo.On("CreateFriend", ctx, "1", "2").Return(nil)

	err = service.CreateFriend(ctx, input)
	assert.Nil(t, err)

	mockRelRepo.ExpectedCalls = nil
	mockRelRepo.On("CheckFriendshipExists", ctx, "1", "2").Return(false, errors.New("database error"))

	err = service.CreateFriend(ctx, input)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "friendship already exists between requestor@example.com and target@example.com")

	mockUserRepo.ExpectedCalls = nil
	mockUserRepo.On("GetUserByEmail", ctx, "requestor@example.com").Return(nil, errors.New("user not found"))

	err = service.CreateFriend(ctx, input)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "user not found")
}

func TestGetFriendListByEmail(t *testing.T) {
	mockRelRepo := new(MockRelationshipRepository)
	mockUserRepo := new(services.MockUserRepository)

	service := services.NewRelationshipService(mockRelRepo, mockUserRepo)

	tests := []struct {
		name          string
		email         string
		setupMocks    func()
		expectedList  []string
		expectedError error
	}{
		{
			name:  "successfully retrieve friend list",
			email: "user@example.com",
			setupMocks: func() {
				mockRelRepo.On("GetFriends", mock.Anything, "user@example.com").
					Return([]string{"friend1@example.com", "friend2@example.com"}, nil)
			},
			expectedList:  []string{"friend1@example.com", "friend2@example.com"},
			expectedError: nil,
		},
		// {
		// 	name:  "error retrieving friend list",
		// 	email: "user@example.com",
		// 	setupMocks: func() {
		// 		mockRelRepo.On("GetFriends", mock.Anything, "user@example.com").
		// 			Return(nil, errors.New("internal server error"))
		// 	},
		// 	expectedList:  nil,
		// 	expectedError: errors.New("internal server error"),
		// },
		// {
		// 	name:  "no friends found",
		// 	email: "user@example.com",
		// 	setupMocks: func() {
		// 		mockRelRepo.On("GetFriends", mock.Anything, "user@example.com").
		// 			Return([]string{}, nil)
		// 	},
		// 	expectedList:  []string{},
		// 	expectedError: nil,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setupMocks()

			friendList, err := service.GetFriendListByEmail(context.Background(), test.email)

			assert.Equal(t, test.expectedList, friendList)
			if test.expectedError != nil {
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRelRepo.AssertExpectations(t)
		})
	}
}
