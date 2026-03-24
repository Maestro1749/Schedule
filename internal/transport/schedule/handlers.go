package schedule

import (
	"schedule/internal/service/schedule"

	"go.uber.org/zap"
)

type ScheduleHandler struct {
	service schedule.ScheduleService
	logger  *zap.Logger
}

func NewScheduleHandler(service schedule.ScheduleService, logger *zap.Logger) *ScheduleHandler {
	return &ScheduleHandler{service: service, logger: logger}
}
