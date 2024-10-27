package relationship

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	response "github.com/koeylp/friends-management/cmd/internal/handler/rest/response"
	"github.com/koeylp/friends-management/cmd/internal/model/dto/relationship/block"
	"github.com/koeylp/friends-management/cmd/internal/model/dto/relationship/friend"
	"github.com/koeylp/friends-management/cmd/internal/model/dto/relationship/subscription"
	"github.com/koeylp/friends-management/cmd/internal/model/dto/user"
	utils "github.com/koeylp/friends-management/cmd/internal/pkg/string_util"
	relationshipRepo "github.com/koeylp/friends-management/cmd/internal/repository/relationship"
	userRepo "github.com/koeylp/friends-management/cmd/internal/repository/user"
)

// RelationshipController defines the interface for relationship-related operations.
type RelationshipController interface {
	CreateFriend(ctx context.Context, friend *friend.CreateFriend) error
	GetFriendListByEmail(ctx context.Context, email string) ([]string, error)
	GetCommonList(ctx context.Context, friend *friend.CommonFriendListReq) ([]string, error)
	Subscribe(ctx context.Context, subscribeReq *subscription.SubscribeRequest) error
	BlockUpdates(ctx context.Context, blockReq *block.BlockRequest) error
	GetUpdatableEmailAddresses(ctx context.Context, recipientReq *subscription.RecipientRequest) ([]string, error)
}

// relationshipControllerImpl implements the RelationshipController interface.
type relationshipControllerImpl struct {
	relationshipRepo relationshipRepo.RelationshipRepository
	userRepo         userRepo.UserRepository
}

// NewRelationshipController creates a new instance of RelationshipController with the provided repositories.
func NewRelationshipController(relationshipRepo relationshipRepo.RelationshipRepository, userRepo userRepo.UserRepository) RelationshipController {
	return &relationshipControllerImpl{relationshipRepo: relationshipRepo, userRepo: userRepo}
}

// CreateFriend handles the creation of a new friendship between two users.
// It checks if a friendship already exists or if there are any blocking updates before creating the friendship.
func (s *relationshipControllerImpl) CreateFriend(ctx context.Context, friend *friend.CreateFriend) error {
	users, err := s.getUsersByEmails(ctx, friend.Friends)
	if err != nil {
		return err
	}
	exists, err := s.relationshipRepo.CheckFriendshipExists(ctx, users[0].ID, users[1].ID)
	if err != nil {
		return fmt.Errorf("failed to check friendship exist: %w", err)
	}

	if exists {
		return response.NewBadRequestError("friendship already exists between " + users[0].Email + " and " + users[1].Email)
	}

	blockExists, err := s.relationshipRepo.CheckBlockExists(ctx, users[0].ID, users[1].ID)
	if err != nil {
		return fmt.Errorf("failed to check blocking updates exist: %w", err)
	}

	if blockExists {
		return response.NewBadRequestError("blocking updates exists between " + users[0].Email + " and " + users[1].Email)
	}

	return s.relationshipRepo.CreateFriend(ctx, users[0].ID, users[1].ID)
}

// GetFriendListByEmail retrieves a list of friends for a user identified by their email.
// It returns a slice of email addresses or an error if retrieval fails.
func (s *relationshipControllerImpl) GetFriendListByEmail(ctx context.Context, email string) ([]string, error) {
	foundUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, response.NewNotFoundError("user not found with email " + email)
	}
	friends, err := s.relationshipRepo.GetFriends(ctx, foundUser.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to retrieve friends: %w", err)
	}

	return friends, nil
}

// getUsersByEmails fetches user details for a list of email addresses.
// It returns a slice of User objects or an error if any user is not found.
func (s *relationshipControllerImpl) getUsersByEmails(ctx context.Context, emails []string) ([]*user.User, error) {
	users := make([]*user.User, len(emails))
	var err error
	for i, email := range emails {
		users[i], err = s.userRepo.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, response.NewBadRequestError("user not found with email " + email)
		}
	}
	return users, nil
}

// GetCommonList retrieves a list of common friends between two users.
// It returns a slice of email addresses or an error if retrieval fails.
func (s *relationshipControllerImpl) GetCommonList(ctx context.Context, friend *friend.CommonFriendListReq) ([]string, error) {
	users, err := s.getUsersByEmails(ctx, friend.Friends)
	if err != nil {
		return nil, err
	}
	commonFriends, err := s.relationshipRepo.GetCommonFriends(ctx, users)
	if err != nil {
		return nil, err
	}
	return commonFriends, err
}

// Subscribe handles the subscription between two users.
// It checks if the requestor and target users exist and if a subscription already exists.
func (s *relationshipControllerImpl) Subscribe(ctx context.Context, subscribeReq *subscription.SubscribeRequest) error {
	requestor, err := s.userRepo.GetUserByEmail(ctx, subscribeReq.Requestor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.NewBadRequestError("requestor not found")
		}
		return fmt.Errorf("failed to retrieve requestor: %w", err)
	}

	target, err := s.userRepo.GetUserByEmail(ctx, subscribeReq.Target)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.NewBadRequestError("target not found")
		}
		return fmt.Errorf("failed to retrieve target: %w", err)
	}

	exists, err := s.relationshipRepo.CheckSubscriptionExists(ctx, requestor.ID, target.ID)
	if err != nil {
		return fmt.Errorf("failed to check subcription exist: %w", err)
	}
	if exists {
		return response.NewBadRequestError("subscription already exists between " + requestor.Email + " and " + target.Email)
	}

	return s.relationshipRepo.Subscribe(ctx, requestor.ID, target.ID)
}

// BlockUpdates handles the request to block updates from a target user.
// It checks if the requestor and target users exist and if a block already exists.
func (s *relationshipControllerImpl) BlockUpdates(ctx context.Context, blockReq *block.BlockRequest) error {
	requestor, err := s.userRepo.GetUserByEmail(ctx, blockReq.Requestor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.NewBadRequestError("requestor not found")
		}
		return fmt.Errorf("failed to retrieve requestor: %w", err)
	}

	target, err := s.userRepo.GetUserByEmail(ctx, blockReq.Target)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.NewBadRequestError("target not found")
		}
		return fmt.Errorf("failed to retrieve target: %w", err)
	}

	exists, err := s.relationshipRepo.CheckBlockExists(ctx, requestor.ID, target.ID)
	if err != nil {
		return fmt.Errorf("failed to check blocking updates exist: %w", err)
	}
	if exists {
		return response.NewBadRequestError("blocking updates already exists between " + requestor.Email + " and " + target.Email)
	}

	return s.relationshipRepo.BlockUpdates(ctx, requestor.ID, target.ID)
}

// GetUpdatableEmailAddresses retrieves email addresses that can be updated based on the sender's context.
// It analyzes mentioned emails in a text and checks if they can be updated.
func (s *relationshipControllerImpl) GetUpdatableEmailAddresses(ctx context.Context, recipientReq *subscription.RecipientRequest) ([]string, error) {
	sender, err := s.userRepo.GetUserByEmail(ctx, recipientReq.Sender)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.NewBadRequestError("sender not found")
		}
		return nil, fmt.Errorf("failed to retrieve requestor: %w", err)
	}

	mentionedEmails := utils.GetEmailFromText(recipientReq.Text)
	users, err := s.getUsersByEmails(ctx, mentionedEmails)
	if err != nil {
		return nil, err
	}

	recipients, err := s.relationshipRepo.GetUpdatableEmailAddresses(ctx, sender.ID)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if !slices.Contains(recipients, user.Email) {
			recipients = append(recipients, user.Email)
		}
	}
	return recipients, nil
}
