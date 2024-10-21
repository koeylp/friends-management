package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/koeylp/friends-management/constants"
	"github.com/koeylp/friends-management/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type RelationshipRepository interface {
	// Friend
	CreateFriend(ctx context.Context, requestor_id, target_id string) error
	CheckFriendshipExists(ctx context.Context, requestor_id, target_id string) (bool, error)
	GetFriends(ctx context.Context, email string) ([]string, error)
	// GetCommonFriends(ctx context.Context, email_1, email_2 string) ([]string, error)
	// IsFriendConnected(ctx context.Context, email_requestor, email_target string) (bool, error)

	// // Block
	// BlockUser(ctx context.Context, email_requestor, email_target string) (bool, error)
	// IsUserBlocked(ctx context.Context, email_requestor, email_target string) (bool, error)

	// // Subcription
	// Subcribe(ctx context.Context, email_requestor, email_target string) (bool, error)

	// GetReceiverUpdatesList(ctx context.Context, email_sender, text string) ([]string, error)

}

type relationshipRepositoryImpl struct {
	db *sql.DB
}

func NewRelationshipRepository(db *sql.DB) RelationshipRepository {
	return &relationshipRepositoryImpl{db: db}
}

func (repo *relationshipRepositoryImpl) CreateFriend(ctx context.Context, requestor_id, target_id string) error {
	friend := models.Relationship{
		ID:               uuid.New().String(),
		RequestorID:      requestor_id,
		TargetID:         target_id,
		RelationshipType: constants.FRIEND,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return friend.Insert(ctx, repo.db, boil.Infer())
}

func (r *relationshipRepositoryImpl) CheckFriendshipExists(ctx context.Context, requestor_id, target_id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM relationships 
		WHERE (requestor_id = $1 AND target_id = $2) OR (requestor_id = $2 AND target_id = $1 AND relationship_type = $3)
	)`
	err := r.db.QueryRowContext(ctx, query, requestor_id, target_id, constants.FRIEND).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (repo *relationshipRepositoryImpl) GetFriends(ctx context.Context, email string) ([]string, error) {
	user, err := models.Users(
		models.UserWhere.Email.EQ(email),
	).One(ctx, repo.db)
	if err != nil {
		return nil, err
	}

	relationships, err := models.Relationships(
		qm.Where("(requestor_id = ? OR target_id = ?) AND relationship_type = ?", user.ID, user.ID, constants.FRIEND),
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
		friend, err := models.Users(
			models.UserWhere.ID.EQ(friendID),
		).One(ctx, repo.db)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend.Email)
	}

	return friends, nil
}
