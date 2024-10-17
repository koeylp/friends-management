package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/koeylp/friends-management/dto/user"
	"github.com/koeylp/friends-management/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *user.CreateUser) error
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (repo *userRepositoryImpl) CreateUser(ctx context.Context, user *user.CreateUser) error {
	friend := models.User{
		ID:        uuid.New().String(),
		Email:     user.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := friend.Insert(ctx, repo.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}
