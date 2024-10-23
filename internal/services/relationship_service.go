package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/koeylp/friends-management/internal/dto/relationship/block"
	"github.com/koeylp/friends-management/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/internal/dto/relationship/subcription"
	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/koeylp/friends-management/internal/repositories"
	"github.com/koeylp/friends-management/internal/responses"
)

type RelationshipService interface {
	CreateFriend(ctx context.Context, friend *friend.CreateFriend) error
	GetFriendListByEmail(ctx context.Context, email string) ([]string, error)
	GetCommonList(ctx context.Context, friend *friend.CommonFriendListReq) ([]string, error)
	Subcribe(ctx context.Context, subscribeReq *subcription.SubscribeRequest) error
	BlockUpdates(ctx context.Context, blockReq *block.BlockRequest) error
	GetUpdatableEmailAddresses(ctx context.Context, recipientReq *subcription.RecipientRequest) ([]string, error)
}

type relationshipServiceImpl struct {
	relationshipRepo repositories.RelationshipRepository
	userRepo         repositories.UserRepository
}

func NewRelationshipService(relationshipRepo repositories.RelationshipRepository, userRepo repositories.UserRepository) RelationshipService {
	return &relationshipServiceImpl{relationshipRepo: relationshipRepo, userRepo: userRepo}
}

func (s *relationshipServiceImpl) CreateFriend(ctx context.Context, friend *friend.CreateFriend) error {
	users, err := s.getUsersByEmails(ctx, friend.Friends)
	if err != nil {
		return err
	}
	exists, err := s.relationshipRepo.CheckFriendshipExists(ctx, users[0].ID, users[1].ID)
	if err != nil {
		return fmt.Errorf("failed to check friendship exist: %w", err)
	}

	if exists {
		return responses.NewBadRequestError("friendship already exists between " + users[0].Email + " and " + users[1].Email)
	}

	blockExists, err := s.relationshipRepo.CheckBlockExists(ctx, users[0].ID, users[1].ID)
	if err != nil {
		return fmt.Errorf("failed to check blocking updates exist: %w", err)
	}

	if blockExists {
		return responses.NewBadRequestError("blocking updates exists between " + users[0].Email + " and " + users[1].Email)
	}

	return s.relationshipRepo.CreateFriend(ctx, users[0].ID, users[1].ID)
}

func (s *relationshipServiceImpl) GetFriendListByEmail(ctx context.Context, email string) ([]string, error) {
	friends, err := s.relationshipRepo.GetFriends(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to retrieve friends: %w", err)
	}

	return friends, nil
}

func (s *relationshipServiceImpl) getUsersByEmails(ctx context.Context, emails []string) ([]*user.User, error) {
	users := make([]*user.User, len(emails))
	var err error
	for i, email := range emails {
		users[i], err = s.userRepo.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, responses.NewBadRequestError("user not found with email " + email)
		}
	}
	return users, nil
}

func (s *relationshipServiceImpl) GetCommonList(ctx context.Context, friend *friend.CommonFriendListReq) ([]string, error) {
	commonFriends, err := s.relationshipRepo.GetCommonFriends(ctx, friend.Friends[0], friend.Friends[1])
	if err != nil {
		return nil, err
	}
	return commonFriends, err
}

func (s *relationshipServiceImpl) Subcribe(ctx context.Context, subscribeReq *subcription.SubscribeRequest) error {
	requestor, err := s.userRepo.GetUserByEmail(ctx, subscribeReq.Requestor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responses.NewBadRequestError("requestor not found")
		}
		return fmt.Errorf("failed to retrieve requestor: %w", err)
	}

	target, err := s.userRepo.GetUserByEmail(ctx, subscribeReq.Target)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responses.NewBadRequestError("target not found")
		}
		return fmt.Errorf("failed to retrieve target: %w", err)
	}

	exists, err := s.relationshipRepo.CheckSubcriptionExists(ctx, requestor.ID, target.ID)
	if err != nil {
		return fmt.Errorf("failed to check subcription exist: %w", err)
	}
	if exists {
		return responses.NewBadRequestError("subscription already exists between " + requestor.Email + " and " + target.Email)
	}

	return s.relationshipRepo.Subcribe(ctx, requestor.ID, target.ID)
}

func (s *relationshipServiceImpl) BlockUpdates(ctx context.Context, blockReq *block.BlockRequest) error {
	requestor, err := s.userRepo.GetUserByEmail(ctx, blockReq.Requestor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responses.NewBadRequestError("requestor not found")
		}
		return fmt.Errorf("failed to retrieve requestor: %w", err)
	}

	target, err := s.userRepo.GetUserByEmail(ctx, blockReq.Target)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responses.NewBadRequestError("target not found")
		}
		return fmt.Errorf("failed to retrieve target: %w", err)
	}

	exists, err := s.relationshipRepo.CheckBlockExists(ctx, requestor.ID, target.ID)
	if err != nil {
		return fmt.Errorf("failed to check locking updates exist: %w", err)
	}
	if exists {
		return responses.NewBadRequestError("blocking updates already exists between " + requestor.Email + " and " + target.Email)
	}

	return s.relationshipRepo.BlockUpdates(ctx, requestor.ID, target.ID)
}

func (s *relationshipServiceImpl) GetUpdatableEmailAddresses(ctx context.Context, recipientReq *subcription.RecipientRequest) ([]string, error) {
	panic("unimplemented")
}
