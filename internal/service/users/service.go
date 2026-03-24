package users

import (
	"schedule/internal/repository/users"

	"go.uber.org/zap"
)

type UserService struct {
	repo   users.UserRepository
	logger *zap.Logger
}

func NewUserService(repo users.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}
