package services

import (
	"context"

	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/koeylp/friends-management/internal/repositories"
)

// UserService defines the interface for user-related operations.
type UserService interface {
	CreateUser(ctx context.Context, user *user.CreateUser) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
}

// userServiceImpl implements the UserService interface.
type userServiceImpl struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new instance of UserService with the provided UserRepository.
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

// CreateUser handles the creation of a new user.
func (s *userServiceImpl) CreateUser(ctx context.Context, user *user.CreateUser) error {
	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return err
}

// GetUserByEmail retrieves a user by their email address.
func (s *userServiceImpl) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
