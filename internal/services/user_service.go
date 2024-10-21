package services

import (
	"context"

	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/koeylp/friends-management/internal/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, user *user.CreateUser) error
}

type userServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user *user.CreateUser) error {
	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return err
}
