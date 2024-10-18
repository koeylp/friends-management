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

type UserRepository interface {
	CreateUser(ctx context.Context, user *user.CreateUser) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

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

func (repo *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	foundUser, err := models.Users(qm.Where("email = ?", email)).One(ctx, repo.db)
	if err != nil {
		return nil, err
	}
	return &user.User{ID: foundUser.ID, Email: foundUser.Email, CreatedAt: foundUser.CreatedAt, UpdatedAt: foundUser.UpdatedAt}, err
}
