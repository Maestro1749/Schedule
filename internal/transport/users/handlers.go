package users

import (
	"schedule/internal/service/users"

	"go.uber.org/zap"
)

type UserHandler struct {
	service users.UserService
	logger  *zap.Logger
}

func NewUserHandler(service users.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}
