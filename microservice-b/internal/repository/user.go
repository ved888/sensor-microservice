package repository

import (
	"database/sql"
	"microservice-b/model"
	"microservice-b/utils"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(user *model.SignupRequest) error {
	query := `INSERT INTO users (
                   first_name, 
                   last_name,
                   email,
                   password,
                   role)
            VALUES (?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.Role)
	return err
}

func (r *UserRepo) GetByEmail(email string) (*model.User, error) {
	u := &model.User{}
	query := `SELECT id,first_name, last_name,email, password, role, archived_at
          FROM users
          WHERE email = ? AND archived_at IS NULL`
	err := r.DB.Get(u, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrEmailNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) UpdateLastLogin(id uint64) error {
	query := `UPDATE users SET 
                        last_login = NOW() WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
