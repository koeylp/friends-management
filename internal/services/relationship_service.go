package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/koeylp/friends-management/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/koeylp/friends-management/internal/repositories"
)

type RelationshipService interface {
	CreateFriend(ctx context.Context, friend *friend.CreateFriend) error
	GetFriendListByEmail(ctx context.Context, email string) ([]string, error)
	GetCommonList(ctx context.Context, friend friend.CommonFriendListReq) ([]string, error)
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
		return fmt.Errorf("friendship already exists between %s and %s", users[0].Email, users[1].Email)
	}

	if exists {
		return fmt.Errorf("friendship already exists between %s and %s", users[0].Email, users[1].Email)
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
			return nil, fmt.Errorf("user not found")
		}
	}
	return users, nil
}

func (s *relationshipServiceImpl) GetCommonList(ctx context.Context, friend friend.CommonFriendListReq) ([]string, error) {
	panic("unimplemented")
}
