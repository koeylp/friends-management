package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/koeylp/friends-management/internal/dto/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	ctx := context.Background()
	userData := &user.CreateUser{
		Email: "test@example.com",
	}

	mock.ExpectExec("INSERT INTO \"users\"").
		WithArgs(sqlmock.AnyArg(), userData.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(ctx, userData)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()
	email := "test@example.com"
	userID := uuid.New().String()
	createdAt := time.Now()
	updatedAt := createdAt

	rows := sqlmock.NewRows([]string{"id", "email", "created_at", "updated_at"}).
		AddRow(userID, email, createdAt, updatedAt)

	mock.ExpectQuery(`SELECT .* FROM "users" WHERE \(email = \$1\) LIMIT 1`).
		WithArgs(email).
		WillReturnRows(rows)

	result, err := repo.GetUserByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, email, result.Email)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, createdAt, result.CreatedAt)
	assert.Equal(t, updatedAt, result.UpdatedAt)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
