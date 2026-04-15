package schedule_service

import (
	"context"
	"schedule/internal/models"
	schedule_repository "schedule/internal/repository/schedule"

	"go.uber.org/zap"
)

type ScheduleService struct {
	repo   schedule_repository.ScheduleRepository
	logger *zap.Logger
}

func NewScheduleService(repo schedule_repository.ScheduleRepository, logger *zap.Logger) *ScheduleService {
	return &ScheduleService{repo: repo, logger: logger}
}

func (s *ScheduleService) GetSchedule(ctx context.Context, req models.GetScheduleRequest) ([]models.ScheduleItemResponse, error) {
	if req.WeekType != 1 && req.WeekType != 2 {
		return nil, models.ErrInvalidWeekType
	}

	if req.Weekday < 1 || req.Weekday > 7 {
		return nil, models.ErrInvalidWeekday
	}

	return s.repo.GetSchedule(ctx, req.GroupID, req.WeekType, req.Weekday, req.Subgroup)
}

func (s *ScheduleService) GetWeekSchedule(ctx context.Context, req models.GetWeekScheduleRequest) ([]models.ScheduleItemResponse, error) {
	if req.WeekType != nil && *req.WeekType != 1 && *req.WeekType != 2 {
		return nil, models.ErrInvalidWeekType
	}

	return s.repo.GetWeekSchedule(ctx, req.GroupID, req.WeekType, req.Subgroup)
}
