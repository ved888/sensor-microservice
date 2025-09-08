package repository_test

import (
	"database/sql"
	"errors"
	"microservice-b/internal/repository"
	"microservice-b/model"
	"microservice-b/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestUserRepo_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repository.NewUserRepository(sqlxDB)

	user := &model.SignupRequest{
		FirstName: "Ved",
		LastName:  "Verma",
		Email:     "ved@example.com",
		Password:  "password123",
		Role:      "admin",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.FirstName, user.LastName, user.Email, user.Password, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(user)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repository.NewUserRepository(sqlxDB)

	// Successful fetch
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "role", "archived_at"}).
		AddRow(1, "Ved", "Verma", "ved@example.com", "password123", "analyst", nil)

	mock.ExpectQuery("SELECT id,first_name, last_name,email, password, role, archived_at").
		WithArgs("ved@example.com").
		WillReturnRows(rows)

	u, err := repo.GetByEmail("ved@example.com")
	require.NoError(t, err)
	require.Equal(t, uint64(1), u.ID)
	require.Equal(t, "Ved", u.FirstName)
	require.Equal(t, "analyst", u.Role)

	// No rows found
	mock.ExpectQuery("SELECT id,first_name, last_name,email, password, role, archived_at").
		WithArgs("notfound@example.com").
		WillReturnError(sql.ErrNoRows)

	u, err = repo.GetByEmail("notfound@example.com")
	require.Nil(t, u)
	require.ErrorIs(t, err, utils.ErrEmailNotFound)

	// Other DB error
	mock.ExpectQuery("SELECT id,first_name, last_name,email, password, role, archived_at").
		WithArgs("error@example.com").
		WillReturnError(errors.New("db failure"))

	u, err = repo.GetByEmail("error@example.com")
	require.Nil(t, u)
	require.EqualError(t, err, "db failure")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_UpdateLastLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repository.NewUserRepository(sqlxDB)

	// Successful update
	mock.ExpectExec("UPDATE users SET").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateLastLogin(1)
	require.NoError(t, err)

	// Update failure
	mock.ExpectExec("UPDATE users SET").
		WithArgs(2).
		WillReturnError(errors.New("db update error"))

	err = repo.UpdateLastLogin(2)
	require.EqualError(t, err, "db update error")

	require.NoError(t, mock.ExpectationsWereMet())
}
