package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/koeylp/friends-management/internal/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// UserRepository defines the interface for user-related database operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user *user.CreateUser) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
}

// userRepositoryImpl implements the UserRepository interface.
type userRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository with the provided database connection.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// CreateUser inserts a new user into the database using the provided user data.
func (repo *userRepositoryImpl) CreateUser(ctx context.Context, user *user.CreateUser) error {
	newUser := models.User{
		ID:        uuid.New().String(),
		Email:     user.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := newUser.Insert(ctx, repo.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

// GetUserByEmail retrieves a user from the database by their email address.
func (repo *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	foundUser, err := models.Users(qm.Where("email = ?", email)).One(ctx, repo.db)
	if err != nil {
		return nil, err
	}
	return &user.User{ID: foundUser.ID, Email: foundUser.Email, CreatedAt: foundUser.CreatedAt, UpdatedAt: foundUser.UpdatedAt}, err
}
