package relationship

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koeylp/friends-management/cmd/internal/model/dto/user"
	"github.com/koeylp/friends-management/cmd/internal/repository/orm"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// RelationshipRepository defines the interface for managing user relationships.
type RelationshipRepository interface {
	// Friend
	CreateFriend(ctx context.Context, requestor_id, target_id string) error
	CheckFriendshipExists(ctx context.Context, requestor_id, target_id string) (bool, error)
	GetFriends(ctx context.Context, email string) ([]string, error)
	GetCommonFriends(ctx context.Context, users []*user.User) ([]string, error)

	// Subscription
	Subscribe(ctx context.Context, requestor_id, target_id string) error
	CheckSubscriptionExists(ctx context.Context, requestor_id, target_id string) (bool, error)
	GetUpdatableEmailAddresses(ctx context.Context, sender_id string) ([]string, error)

	// Block
	BlockUpdates(ctx context.Context, requestor_id, target_id string) error
	CheckBlockExists(ctx context.Context, requestor_id, target_id string) (bool, error)
}

// relationshipRepositoryImpl is the implementation of the RelationshipRepository interface.
type relationshipRepositoryImpl struct {
	db *sql.DB
}

// NewRelationshipRepository creates a new instance of RelationshipRepository.
func NewRelationshipRepository(db *sql.DB) RelationshipRepository {
	return &relationshipRepositoryImpl{db: db}
}

// CreateFriend adds a new friendship relationship to the database.
func (repo *relationshipRepositoryImpl) CreateFriend(ctx context.Context, requestor_id, target_id string) error {
	friend := orm.Relationship{
		ID:               uuid.New().String(),
		RequestorID:      requestor_id,
		TargetID:         target_id,
		RelationshipType: FRIEND,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return friend.Insert(ctx, repo.db, boil.Infer())
}

// CheckFriendshipExists checks if a friendship exists between two users.
func (repo *relationshipRepositoryImpl) CheckFriendshipExists(ctx context.Context, requestor_id, target_id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM relationships 
		WHERE (requestor_id = $1 AND target_id = $2 AND relationship_type = $3) OR (requestor_id = $2 AND target_id = $1 AND relationship_type = $3)
	)`
	err := repo.db.QueryRowContext(ctx, query, requestor_id, target_id, FRIEND).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetFriends retrieves the list of friends for a given user by email.
func (repo *relationshipRepositoryImpl) GetFriends(ctx context.Context, email string) ([]string, error) {
	user, err := orm.Users(
		orm.UserWhere.Email.EQ(email),
	).One(ctx, repo.db)
	if err != nil {
		return nil, err
	}

	relationships, err := orm.Relationships(
		qm.Where("(requestor_id = ? OR target_id = ?) AND relationship_type = ?", user.ID, user.ID, FRIEND),
	).All(ctx, repo.db)
	if err != nil {
		return nil, err
	}

	friendIDs := make([]string, 0)
	for _, relationship := range relationships {
		if relationship.RequestorID == user.ID {
			friendIDs = append(friendIDs, relationship.TargetID)
		} else {
			friendIDs = append(friendIDs, relationship.RequestorID)
		}
	}

	friends := make([]string, 0)
	for _, friendID := range friendIDs {
		friend, err := orm.Users(
			orm.UserWhere.ID.EQ(friendID),
		).One(ctx, repo.db)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend.Email)
	}

	return friends, nil
}

// GetCommonFriends retrieves the common friends between two users.
func (repo *relationshipRepositoryImpl) GetCommonFriends(ctx context.Context, users []*user.User) ([]string, error) {
	query := `
    WITH user1_friends AS (
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
    JOIN user2_friends u2
        ON u1.friend_id = u2.friend_id
	JOIN users u
    ON u1.friend_id = u.id;
    `

	rows, err := repo.db.QueryContext(ctx, query, users[0].ID, users[1].ID, FRIEND)
	if err != nil {
		return nil, fmt.Errorf("failed to query common friends: %w", err)
	}
	defer rows.Close()

	var commonFriends []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, fmt.Errorf("failed to scan email: %w", err)
		}
		commonFriends = append(commonFriends, email)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return commonFriends, nil
}

// Subscribe adds a new subscription relationship to the database.
func (repo *relationshipRepositoryImpl) Subscribe(ctx context.Context, requestor_id string, target_id string) error {
	subcription := orm.Relationship{
		ID:               uuid.New().String(),
		RequestorID:      requestor_id,
		TargetID:         target_id,
		RelationshipType: SUBSCRIBE,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return subcription.Insert(ctx, repo.db, boil.Infer())
}

// CheckSubscriptionExists checks if a subscription relationship exists between two users.
func (repo *relationshipRepositoryImpl) CheckSubscriptionExists(ctx context.Context, requestor_id string, target_id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM relationships 
		WHERE (requestor_id = $1 AND target_id = $2 AND relationship_type = $3) OR (requestor_id = $2 AND target_id = $1 AND relationship_type = $3)
	)`
	err := repo.db.QueryRowContext(ctx, query, requestor_id, target_id, SUBSCRIBE).Scan(&exists)
	if err != nil {
		return true, err
	}
	return exists, nil
}

// CheckBlockExists checks if a block relationship exists between two users.
func (repo *relationshipRepositoryImpl) CheckBlockExists(ctx context.Context, requestor_id string, target_id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM relationships 
		WHERE (requestor_id = $1 AND target_id = $2 AND relationship_type = $3) OR (requestor_id = $2 AND target_id = $1 AND relationship_type = $3)
	)`
	err := repo.db.QueryRowContext(ctx, query, requestor_id, target_id, BLOCK).Scan(&exists)
	if err != nil {
		return true, err
	}
	return exists, nil
}

// BlockUpdates adds a new block relationship to the database.
func (repo *relationshipRepositoryImpl) BlockUpdates(ctx context.Context, requestor_id string, target_id string) error {
	block := orm.Relationship{
		ID:               uuid.New().String(),
		RequestorID:      requestor_id,
		TargetID:         target_id,
		RelationshipType: BLOCK,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return block.Insert(ctx, repo.db, boil.Infer())
}

// GetUpdatableEmailAddresses retrieves email addresses that can be updated, filtering out blocked users.
func (repo *relationshipRepositoryImpl) GetUpdatableEmailAddresses(ctx context.Context, sender_id string) ([]string, error) {
	recipients, err := orm.Users(
		qm.Select("users.email"),
		qm.Distinct("users.email"),
		qm.Where("users.id != ?", sender_id),
		qm.LeftOuterJoin("relationships AS r1 ON (r1.requestor_id = users.id AND r1.target_id = ?) OR (r1.target_id = users.id AND r1.requestor_id = ?)", sender_id, sender_id),
		qm.LeftOuterJoin("relationships AS r2 ON r2.requestor_id = users.id AND r2.target_id = ? AND r2.relationship_type = ?", sender_id, SUBSCRIBE),
		qm.Where("users.id NOT IN (SELECT target_id FROM relationships WHERE requestor_id = ? AND relationship_type = ?)", sender_id, BLOCK),
		qm.Where("r1.relationship_type = 'Friend' OR r2.relationship_type = 'Subscribe'"),
	).All(ctx, repo.db)

	if err != nil {
		return nil, fmt.Errorf("failed to get recipients: %v", err)
	}

	emails := make([]string, 0, len(recipients))
	for _, user := range recipients {
		emails = append(emails, user.Email)
	}
	return emails, nil
}
