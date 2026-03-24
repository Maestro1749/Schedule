package users

import (
	"database/sql"

	"go.uber.org/zap"
)

type UserRepository interface {
}

type userRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserRepository(db *sql.DB, logger *zap.Logger) userRepo {
	return userRepo{db: db, logger: logger}
}
