package services

import (
	"context"

	"github.com/koeylp/friends-management/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/koeylp/friends-management/internal/repository"
)

type RelationshipService struct {
	relationshipRepo repository.RelationshipRepository
	userRepo         repository.UserRepository
}

func NewRelationshipService(relationshipRepo repository.RelationshipRepository, userRepo repository.UserRepository) *RelationshipService {
	return &RelationshipService{relationshipRepo: relationshipRepo, userRepo: userRepo}
}

func (s *RelationshipService) CreateFriend(ctx context.Context, friend *friend.CreateFriend) error {
	users, err := s.getUsersByEmails(ctx, friend.Friends)
	if err != nil {
		return err
	}

	return s.relationshipRepo.CreateFriend(ctx, users[0].ID, users[1].ID)
}

func (s *RelationshipService) getUsersByEmails(ctx context.Context, emails []string) ([]*user.User, error) {
	users := make([]*user.User, len(emails))
	var err error
	for i, email := range emails {
		users[i], err = s.userRepo.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}
