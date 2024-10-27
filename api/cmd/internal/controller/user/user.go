package user

import (
	"context"

	"github.com/koeylp/friends-management/cmd/internal/model/dto/user"
	userRepo "github.com/koeylp/friends-management/cmd/internal/repository/user"
)

// UserController defines the interface for user-related operations.
type UserController interface {
	CreateUser(ctx context.Context, user *user.CreateUser) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
}

// userControllerImpl implements the UserController interface.
type userControllerImpl struct {
	userRepo userRepo.UserRepository
}

// NewUserController creates a new instance of UserController with the provided UserRepository.
func NewUserController(userRepo userRepo.UserRepository) UserController {
	return &userControllerImpl{userRepo: userRepo}
}

// CreateUser handles the creation of a new user.
func (s *userControllerImpl) CreateUser(ctx context.Context, user *user.CreateUser) error {
	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return err
}

// GetUserByEmail retrieves a user by their email address.
func (s *userControllerImpl) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
