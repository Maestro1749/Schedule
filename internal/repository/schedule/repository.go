package schedule

import (
	"database/sql"

	"go.uber.org/zap"
)

type ScheduleRepository interface {
}

type scheduleRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewScheduleRepo(db *sql.DB, logger *zap.Logger) scheduleRepo {
	return scheduleRepo{db: db, logger: logger}
}
