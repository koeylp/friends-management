package services_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/koeylp/friends-management/internal/dto/relationship/block"
	"github.com/koeylp/friends-management/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/internal/dto/relationship/subcription"
	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/koeylp/friends-management/internal/repositories"
	"github.com/koeylp/friends-management/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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

func (m *MockRelationshipRepository) CheckSubcriptionExists(ctx context.Context, requestor_id string, target_id string) (bool, error) {
	args := m.Called(ctx, requestor_id, target_id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRelationshipRepository) Subcribe(ctx context.Context, requestor_id string, target_id string) error {
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

func TestRelationshipService_GetCommonList(t *testing.T) {
	ctx := context.Background()
	mockRelRepo := new(MockRelationshipRepository)
	mockUserRepo := new(services.MockUserRepository)
	mockService := services.NewRelationshipService(mockRelRepo, mockUserRepo)

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
	mockUserRepo := new(services.MockUserRepository)

	service := services.NewRelationshipService(mockRelRepo, mockUserRepo)

	requestor := &user.User{ID: "123", Email: "requestor@example.com"}
	target := &user.User{ID: "456", Email: "target@example.com"}

	subscribeReq := &subcription.SubscribeRequest{
		Requestor: "requestor@example.com",
		Target:    "target@example.com",
	}

	mockUserRepo.On("GetUserByEmail", ctx, "requestor@example.com").Return(requestor, nil)
	mockUserRepo.On("GetUserByEmail", ctx, "target@example.com").Return(target, nil)

	mockRelRepo.On("CheckSubcriptionExists", ctx, requestor.ID, target.ID).Return(false, nil)

	mockRelRepo.On("Subcribe", ctx, requestor.ID, target.ID).Return(nil)

	err := service.Subcribe(ctx, subscribeReq)

	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockRelRepo.AssertExpectations(t)
}

func TestBlockUpdates(t *testing.T) {
	ctx := context.Background()

	mockRelRepo := new(MockRelationshipRepository)
	mockUserRepo := new(services.MockUserRepository)

	service := services.NewRelationshipService(mockRelRepo, mockUserRepo)

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
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewRelationshipRepository(db)

	ctx := context.Background()
	senderID := "123"
	mentionedEmails := []string{"email1@example.com", "email2@example.com"} // Example non-empty slice

	// Create placeholders for the IN clause
	placeholders := strings.Repeat("?,", len(mentionedEmails)-1) + "?"

	// SQL query with placeholders
	query := `
		SELECT DISTINCT users.email
		FROM users
		LEFT JOIN relationships AS r1
			ON (r1.requestor_id = users.id AND r1.target_id = $1)
			OR (r1.target_id = users.id AND r1.requestor_id = $2)
		LEFT JOIN relationships AS r2
			ON r2.requestor_id = users.id
			AND r2.target_id = $3
			AND r2.relationship_type = $4
		WHERE users.id != $5
		AND users.id NOT IN (
			SELECT target_id FROM relationships WHERE requestor_id = $6 AND relationship_type = $7
		)
		AND users.email IN (` + placeholders + `)`

	// Mock results
	rows := sqlmock.NewRows([]string{"email"}).
		AddRow("email1@example.com").
		AddRow("email3@example.com")

	// Prepare arguments
	args := []driver.Value{
		senderID, senderID, senderID, "Subscribe", senderID, senderID, "Block",
	}
	for _, email := range mentionedEmails {
		args = append(args, email)
	}

	// Expectation for the query
	mock.ExpectQuery(query).
		WithArgs(args...).
		WillReturnRows(rows)

	// Call the method
	emails, err := repo.GetUpdatableEmailAddresses(ctx, senderID)

	// Validate the result
	assert.NoError(t, err)
	assert.Equal(t, []string{"email1@example.com", "email3@example.com"}, emails)

	// Mock an error case
	mock.ExpectQuery(query).
		WithArgs(args...).
		WillReturnError(errors.New("db error"))

	// Call again, expecting an error
	emails, err = repo.GetUpdatableEmailAddresses(ctx, senderID)

	// Validate error case
	assert.Error(t, err)
	assert.Nil(t, emails)

	// Ensure all expectations are met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
