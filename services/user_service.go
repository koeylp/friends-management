package services

import (
	"context"

	"github.com/koeylp/friends-management/dto/user"
	"github.com/koeylp/friends-management/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(friendRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: friendRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user *user.CreateUser) error {
	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return err
}
