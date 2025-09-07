package model

import "time"

type User struct {
	ID         uint64     `db:"id" json:"id"`
	FirstName  string     `db:"first_name" json:"first_name"`
	LastName   string     `db:"last_name" json:"last_name"`
	Email      string     `db:"email" json:"email"`
	Password   string     `db:"password" json:"password"`
	Role       string     `db:"role" json:"role"`
	LastLogin  *time.Time `db:"last_login" json:"last_login"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	ArchivedAt *time.Time `db:"archived_at" json:"archived_at,omitempty"`
}

type SignupRequest struct {
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"password"`
	Role      string `db:"role" json:"role"`
}
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type SignupResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}
