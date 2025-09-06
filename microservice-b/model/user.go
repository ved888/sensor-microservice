package model

import "time"

type User struct {
	ID         uint64     `db:"id" json:"id"`
	Email      string     `db:"email" json:"email"`
	Password   string     `db:"password" json:"-"`
	Role       string     `db:"role" json:"role"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	ArchivedAt *time.Time `db:"archived_at" json:"archived_at,omitempty"`
}
