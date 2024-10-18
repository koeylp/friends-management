package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/koeylp/friends-management/constants"
	"github.com/koeylp/friends-management/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type RelationshipRepository interface {
	// Friend
	CreateFriend(ctx context.Context, requestor_id, target_id string) error
	// GetFriends(ctx context.Context, email string) ([]string, error)
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
