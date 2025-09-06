package config

import (
	"github.com/jmoiron/sqlx"
)

type DAO struct {
	DB *sqlx.DB
}
