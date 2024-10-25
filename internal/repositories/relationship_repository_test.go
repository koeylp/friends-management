package repositories_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/koeylp/friends-management/constants"
	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/koeylp/friends-management/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFriend(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewRelationshipRepository(db)

	requestorID := "user1-id"
	targetID := "user2-id"

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "relationships" ("id","requestor_id","target_id","relationship_type","created_at","updated_at")`)).
		WithArgs(
			sqlmock.AnyArg(),
			requestorID,
			targetID,
			constants.FRIEND,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateFriend(context.Background(), requestorID, targetID)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCheckFriendshipExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewRelationshipRepository(db)

	requestorID := "123"
	targetID := "456"
	relationshipType := constants.FRIEND

	// Test Case: Friendship exists
	mock.ExpectQuery(`SELECT EXISTS \(.*\)`).
		WithArgs(requestorID, targetID, relationshipType).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.CheckFriendshipExists(context.Background(), requestorID, targetID)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test Case: Friendship does not exist
	mock.ExpectQuery(`SELECT EXISTS \(.*\)`).
		WithArgs(requestorID, targetID, relationshipType).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err = repo.CheckFriendshipExists(context.Background(), requestorID, targetID)
	assert.NoError(t, err)
	assert.False(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFriends(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewRelationshipRepository(db)

	email := "test@example.com"
	userID := "123"
	friendID := "456"
	friendEmail := "friend@example.com"

	// Mock the query to fetch the user by email
	mock.ExpectQuery(`SELECT "users"\.\* FROM "users" WHERE \("users"\."email" = \$1\) LIMIT 1`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(userID, email))

	// Mock the query to fetch the relationships for the user
	mock.ExpectQuery(`SELECT "relationships"\.\* FROM "relationships" WHERE \(\(requestor_id = \$1 OR target_id = \$2\) AND relationship_type = \$3\)`).
		WithArgs(userID, userID, constants.FRIEND).
		WillReturnRows(sqlmock.NewRows([]string{"requestor_id", "target_id"}).AddRow(userID, friendID))

	// Mock the query to fetch the friend by ID
	mock.ExpectQuery(`SELECT "users"\.\* FROM "users" WHERE \("users"\."id" = \$1\) LIMIT 1`).
		WithArgs(friendID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(friendID, friendEmail))

	friends, err := repo.GetFriends(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, []string{friendEmail}, friends)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCommonFriends(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewRelationshipRepository(db)

	users := []*user.User{
		{ID: "user1-id", Email: "user1@example.com"},
		{ID: "user2-id", Email: "user2@example.com"},
	}

	query := `WITH user1_friends AS (
	SELECT CASE
		WHEN r1.requestor_id = $1 THEN r1.target_id
		ELSE r1.requestor_id
	END AS friend_id
	FROM relationships r1
	WHERE (r1.requestor_id = $1 OR r1.target_id = $1)
	  AND r1.relationship_type = $3
  ),
  user2_friends AS (
	SELECT CASE
		WHEN r2.requestor_id = $2 THEN r2.target_id
		ELSE r2.requestor_id
	END AS friend_id
	FROM relationships r2
	WHERE (r2.requestor_id = $2 OR r2.target_id = $2)
	  AND r2.relationship_type = $3
  )
  SELECT DISTINCT u.email
  FROM user1_friends u1
  JOIN user2_friends u2 ON u1.friend_id = u2.friend_id
  JOIN users u ON u1.friend_id = u.id;`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(users[0].ID, users[1].ID, constants.FRIEND).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
			AddRow("commonFriend1@example.com").
			AddRow("commonFriend2@example.com"))

	result, err := repo.GetCommonFriends(context.Background(), users)

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"commonFriend1@example.com", "commonFriend2@example.com"}, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSubscribe(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewRelationshipRepository(db)

	requestorID := "user1-id"
	targetID := "user2-id"

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "relationships" ("id","requestor_id","target_id","relationship_type","created_at","updated_at")`)).
		WithArgs(
			sqlmock.AnyArg(),
			requestorID,
			targetID,
			constants.SUBSCRIBE,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Subscribe(context.Background(), requestorID, targetID)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCheckSubscriptionExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewRelationshipRepository(db)

	requestorID := "123"
	targetID := "456"
	relationshipType := constants.SUBSCRIBE

	// Test Case: Subscription exists
	mock.ExpectQuery(`SELECT EXISTS \(.*\)`).
		WithArgs(requestorID, targetID, relationshipType).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.CheckSubscriptionExists(context.Background(), requestorID, targetID)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test Case: Subscription does not exist
	mock.ExpectQuery(`SELECT EXISTS \(.*\)`).
		WithArgs(requestorID, targetID, relationshipType).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err = repo.CheckSubscriptionExists(context.Background(), requestorID, targetID)
	assert.NoError(t, err)
	assert.False(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckBlockExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewRelationshipRepository(db)

	requestorID := "123"
	targetID := "456"
	relationshipType := constants.BLOCK

	// Test Case: Block exists
	mock.ExpectQuery(`SELECT EXISTS \(.*\)`).
		WithArgs(requestorID, targetID, relationshipType).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.CheckBlockExists(context.Background(), requestorID, targetID)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test Case: Block does not exist
	mock.ExpectQuery(`SELECT EXISTS \(.*\)`).
		WithArgs(requestorID, targetID, relationshipType).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err = repo.CheckBlockExists(context.Background(), requestorID, targetID)
	assert.NoError(t, err)
	assert.False(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewRelationshipRepository(db)

	requestorID := "user1-id"
	targetID := "user2-id"

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "relationships" ("id","requestor_id","target_id","relationship_type","created_at","updated_at")`)).
		WithArgs(
			sqlmock.AnyArg(),
			requestorID,
			targetID,
			constants.BLOCK,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.BlockUpdates(context.Background(), requestorID, targetID)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUpdatableEmailAddresses(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repositories.NewRelationshipRepository(db)

	ctx := context.Background()
	senderID := "123"

	rows := sqlmock.NewRows([]string{"email"}).
		AddRow("test1@example.com").
		AddRow("friend@example.com")

	mock.ExpectQuery(`SELECT DISTINCT users.email FROM "users" LEFT JOIN relationships AS r1 ON \(r1\.requestor_id = users\.id AND r1\.target_id = \$1\) OR \(r1\.target_id = users\.id AND r1\.requestor_id = \$2\) LEFT JOIN relationships AS r2 ON r2\.requestor_id = users\.id AND r2\.target_id = \$3 AND r2\.relationship_type = \$4 WHERE \(users\.id != \$5\) AND \(users\.id NOT IN \(SELECT target_id FROM relationships WHERE requestor_id = \$6 AND relationship_type = \$7\)\) AND \(r1\.relationship_type = 'Friend' OR r2\.relationship_type = 'Subscribe'\)`).
		WithArgs(senderID, senderID, senderID, constants.SUBSCRIBE, senderID, senderID, constants.BLOCK).
		WillReturnRows(rows)

	emails, err := repo.GetUpdatableEmailAddresses(ctx, senderID)

	require.NoError(t, err)
	require.Equal(t, []string{"test1@example.com", "friend@example.com"}, emails)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
