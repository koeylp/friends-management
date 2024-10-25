package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/koeylp/friends-management/internal/dto/relationship/block"
	"github.com/koeylp/friends-management/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/internal/dto/relationship/subscription"
	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateFriend(t *testing.T) {
	ctx := context.Background()

	mockRelRepo := new(MockRelationshipRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewRelationshipService(mockRelRepo, mockUserRepo)

	inputEmails := []string{"requestor@example.com", "target@example.com"}
	input := &friend.CreateFriend{
		Friends: inputEmails,
	}

	mockUsers := []*user.User{
		{ID: "1", Email: "requestor@example.com"},
		{ID: "2", Email: "target@example.com"},
	}

	// Case 1: Friendship already exists
	mockUserRepo.On("GetUserByEmail", ctx, "requestor@example.com").Return(mockUsers[0], nil)
	mockUserRepo.On("GetUserByEmail", ctx, "target@example.com").Return(mockUsers[1], nil)
	mockRelRepo.On("CheckFriendshipExists", ctx, "1", "2").Return(true, nil)

	err := service.CreateFriend(ctx, input)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "400: friendship already exists between requestor@example.com and target@example.com")

	// Reset expected calls
	mockRelRepo.ExpectedCalls = nil

	// Case 2: Successful friend creation (no block)
	// mockRelRepo.On("CheckFriendshipExists", ctx, "1", "2").Return(false, nil)
	// mockRelRepo.On("CheckBlockExists", ctx, "1", "2").Return(false, nil)
	// mockRelRepo.On("CreateFriend", ctx, "1", "2").Return(nil)

	// err = service.CreateFriend(ctx, input)
	// assert.Nil(t, err)

	// // Reset expected calls
	// mockRelRepo.ExpectedCalls = nil

	// Case 3: Block exists between the users
	mockRelRepo.On("CheckFriendshipExists", ctx, "1", "2").Return(false, nil)
	mockRelRepo.On("CheckBlockExists", ctx, "1", "2").Return(true, nil)

	err = service.CreateFriend(ctx, input)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "400: blocking updates exists between requestor@example.com and target@example.com")

	// Reset expected calls
	mockRelRepo.ExpectedCalls = nil

	// Case 4: Error while checking block existence
	mockRelRepo.On("CheckFriendshipExists", ctx, "1", "2").Return(false, nil)
	mockRelRepo.On("CheckBlockExists", ctx, "1", "2").Return(false, errors.New("database error"))

	err = service.CreateFriend(ctx, input)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "failed to check blocking updates exist: database error")

	// Reset expected calls
	mockRelRepo.ExpectedCalls = nil

	// Case 5: Error while checking friendship existence
	mockRelRepo.On("CheckFriendshipExists", ctx, "1", "2").Return(false, errors.New("database error"))

	err = service.CreateFriend(ctx, input)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "failed to check friendship exist: database error")

	// Reset expected calls for user repository
	mockUserRepo.ExpectedCalls = nil

	// Case 6: User not found (requestor)
	mockUserRepo.On("GetUserByEmail", ctx, "requestor@example.com").Return(nil, errors.New("user not found"))

	err = service.CreateFriend(ctx, input)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "400: user not found with email requestor@example.com")

	// Reset expected calls for user repository
	mockUserRepo.ExpectedCalls = nil

	// Case 7: User not found (target)
	mockUserRepo.On("GetUserByEmail", ctx, "requestor@example.com").Return(mockUsers[0], nil)
	mockUserRepo.On("GetUserByEmail", ctx, "target@example.com").Return(nil, errors.New("user not found"))

	err = service.CreateFriend(ctx, input)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "400: user not found with email target@example.com")
}

func TestGetFriendListByEmail(t *testing.T) {
	mockRelRepo := new(MockRelationshipRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewRelationshipService(mockRelRepo, mockUserRepo)

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

func TestRelationshipService_GetCommonList(t *testing.T) {
	ctx := context.Background()
	mockRelRepo := new(MockRelationshipRepository)
	mockUserRepo := new(MockUserRepository)
	mockService := NewRelationshipService(mockRelRepo, mockUserRepo)

	req := &friend.CommonFriendListReq{
		Friends: []string{"user@example.com", "user1@example.com"},
	}

	users := []*user.User{
		{ID: "1", Email: "user@example.com"},
		{ID: "2", Email: "user1@example.com"},
	}
	mockUserRepo.On("GetUserByEmail", ctx, "user@example.com").Return(users[0], nil)
	mockUserRepo.On("GetUserByEmail", ctx, "user1@example.com").Return(users[1], nil)

	expectedCommonFriends := []string{"common.friend1@example.com", "common.friend2@example.com"}

	mockRelRepo.On("GetCommonFriends", ctx, users).Return(expectedCommonFriends, nil)

	commonFriends, err := mockService.GetCommonList(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedCommonFriends, commonFriends)

	mockRelRepo.AssertExpectations(t)
}

func TestSubcribe_Success(t *testing.T) {
	ctx := context.Background()

	mockRelRepo := new(MockRelationshipRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewRelationshipService(mockRelRepo, mockUserRepo)

	requestor := &user.User{ID: "123", Email: "requestor@example.com"}
	target := &user.User{ID: "456", Email: "target@example.com"}

	subscribeReq := &subscription.SubscribeRequest{
		Requestor: "requestor@example.com",
		Target:    "target@example.com",
	}

	mockUserRepo.On("GetUserByEmail", ctx, "requestor@example.com").Return(requestor, nil)
	mockUserRepo.On("GetUserByEmail", ctx, "target@example.com").Return(target, nil)

	mockRelRepo.On("CheckSubscriptionExists", ctx, requestor.ID, target.ID).Return(false, nil)

	mockRelRepo.On("Subscribe", ctx, requestor.ID, target.ID).Return(nil)

	err := service.Subscribe(ctx, subscribeReq)

	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockRelRepo.AssertExpectations(t)
}

func TestBlockUpdates(t *testing.T) {
	ctx := context.Background()

	mockRelRepo := new(MockRelationshipRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewRelationshipService(mockRelRepo, mockUserRepo)

	inputEmails := &block.BlockRequest{
		Requestor: "requestor@example.com",
		Target:    "target@example.com",
	}

	mockUsers := []*user.User{
		{ID: "1", Email: "requestor@example.com"},
		{ID: "2", Email: "target@example.com"},
	}

	// Case 1: User not found (requestor)
	mockUserRepo.On("GetUserByEmail", ctx, "requestor@example.com").Return(nil, errors.New("user not found"))
	err := service.BlockUpdates(ctx, inputEmails)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "failed to retrieve requestor: user not found")

	// Reset expected calls
	mockUserRepo.ExpectedCalls = nil

	// Case 2: User not found (target)
	mockUserRepo.On("GetUserByEmail", ctx, "requestor@example.com").Return(mockUsers[0], nil)
	mockUserRepo.On("GetUserByEmail", ctx, "target@example.com").Return(nil, errors.New("user not found"))
	err = service.BlockUpdates(ctx, inputEmails)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "failed to retrieve target: user not found")

	// Reset expected calls
	mockUserRepo.ExpectedCalls = nil

	// Case 3: Block relationship already exists
	mockUserRepo.On("GetUserByEmail", ctx, "requestor@example.com").Return(mockUsers[0], nil)
	mockUserRepo.On("GetUserByEmail", ctx, "target@example.com").Return(mockUsers[1], nil)
	mockRelRepo.On("CheckBlockExists", ctx, "1", "2").Return(true, nil)
	err = service.BlockUpdates(ctx, inputEmails)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "400: blocking updates already exists between requestor@example.com and target@example.com")

	// Reset expected calls
	mockRelRepo.ExpectedCalls = nil

	// Case 4: Successful block
	mockRelRepo.On("CheckBlockExists", ctx, "1", "2").Return(false, nil)
	mockRelRepo.On("BlockUpdates", ctx, "1", "2").Return(nil)
	err = service.BlockUpdates(ctx, inputEmails)
	assert.Nil(t, err)

	// Case 5: Error while checking block existence
	mockRelRepo.ExpectedCalls = nil
	mockRelRepo.On("CheckBlockExists", ctx, "1", "2").Return(false, errors.New("database error"))
	err = service.BlockUpdates(ctx, inputEmails)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "failed to check blocking updates exist: database error")
}

func TestGetUpdatableEmailAddresses(t *testing.T) {
	ctx := context.Background()

	mockRelRepo := new(MockRelationshipRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewRelationshipService(mockRelRepo, mockUserRepo)

	recipientReq := &subscription.RecipientRequest{
		Sender: "sender@example.com",
		Text:   "This is a test email with some@example.com",
	}
	sender := &user.User{ID: "1", Email: "sender@example.com"}
	userMentioned := &user.User{ID: "2", Email: "some@example.com"}
	updatableEmails := []string{"existing@example.com"}

	// Case 1: Sender not found
	mockUserRepo.On("GetUserByEmail", ctx, "sender@example.com").Return(nil, sql.ErrNoRows)
	recipients, err := service.GetUpdatableEmailAddresses(ctx, recipientReq)
	assert.Nil(t, recipients)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "400: sender not found")

	mockUserRepo.ExpectedCalls = nil

	// Case 2: Successful retrieval of updatable email addresses
	mockUserRepo.On("GetUserByEmail", ctx, "sender@example.com").Return(sender, nil)
	mockUserRepo.On("GetUserByEmail", ctx, "some@example.com").Return(userMentioned, nil)
	mockRelRepo.On("GetUpdatableEmailAddresses", ctx, sender.ID).Return(updatableEmails, nil)

	recipients, err = service.GetUpdatableEmailAddresses(ctx, recipientReq)
	assert.Nil(t, err)
	assert.Contains(t, recipients, "existing@example.com")
	assert.Contains(t, recipients, "some@example.com")

	mockUserRepo.ExpectedCalls = nil
	mockRelRepo.ExpectedCalls = nil

	// Case 3: Error in retrieving updatable email addresses
	mockUserRepo.On("GetUserByEmail", ctx, "sender@example.com").Return(sender, nil)
	mockUserRepo.On("GetUserByEmail", ctx, "some@example.com").Return(userMentioned, nil)
	mockRelRepo.On("GetUpdatableEmailAddresses", ctx, sender.ID).Return([]string(nil), errors.New("db error"))

	recipients, err = service.GetUpdatableEmailAddresses(ctx, recipientReq)
	assert.Nil(t, recipients)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "db error")
}
