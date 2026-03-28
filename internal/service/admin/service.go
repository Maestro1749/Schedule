package admin_service

import (
	"schedule/internal/models"
	admin_repository "schedule/internal/repository/admin"

	"go.uber.org/zap"
)

type AdminService struct {
	repo   admin_repository.AdminRepository
	logger *zap.Logger
}

func NewUserService(repo admin_repository.AdminRepository, logger *zap.Logger) *AdminService {
	return &AdminService{repo: repo, logger: logger}
}

func (s *AdminService) CreateSchedule(data []models.CreateScheduleDTO) error {
	if len(data) == 0 {
		return models.ErrInvalidDataInput
	}

	for _, item := range data {
		if item.GroupID == 0 || item.SubjectID == 0 || item.TeacherID == 0 || item.ClassroomID == 0 {
			return models.ErrInvalidDataInput
		}
		if item.Weekday < 1 || item.Weekday > 7 {
			return models.ErrInvalidWeekday
		}
		if item.WeekType != 1 && item.WeekType != 2 {
			return models.ErrInvalidWeekType
		}
	}

	if err := s.repo.CreateSchedule(data); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) AddTeacher(teachers []models.Teacher) error {
	if len(teachers) == 0 {
		return models.ErrInvalidDataInput
	}

	for _, t := range teachers {
		if t.Fullname == "" {
			return models.ErrInvalidDataInput
		}
	}

	if err := s.repo.AddTeacher(teachers); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) AddClassroom(classrooms []models.Classroom) error {
	if len(classrooms) == 0 {
		return models.ErrInvalidDataInput
	}

	for _, c := range classrooms {
		if c.Number == "" {
			return models.ErrInvalidDataInput
		}
	}

	if err := s.repo.AddClassroom(classrooms); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) AddSubject(subjects []models.Subject) error {
	if len(subjects) == 0 {
		return models.ErrInvalidDataInput
	}

	for _, s := range subjects {
		if s.Name == "" {
			return models.ErrInvalidDataInput
		}
	}

	if err := s.repo.AddSubject(subjects); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) AddGroup(groups []models.Group) error {
	if len(groups) == 0 {
		return models.ErrInvalidDataInput
	}

	for _, g := range groups {
		if g.Name == "" {
			return models.ErrInvalidDataInput
		}
	}

	if err := s.repo.AddGroup(groups); err != nil {
		return models.ErrInvalidDataInput
	}

	return nil
}
