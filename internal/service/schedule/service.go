package schedule

import (
	"schedule/internal/repository/schedule"

	"go.uber.org/zap"
)

type ScheduleService struct {
	repo   schedule.ScheduleRepository
	logger *zap.Logger
}

func NewScheduleService(repo schedule.ScheduleRepository, logger *zap.Logger) *ScheduleService {
	return &ScheduleService{repo: repo, logger: logger}
}
