package admin_service

import (
	"schedule/internal/models"
	admin_repository "schedule/internal/repository/admin"
	"strings"

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
		if item.WeekType != nil && *item.WeekType != 1 && *item.WeekType != 2 {
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

	seen := make(map[string]struct{})

	for i := range teachers {
		teachers[i].Fullname = strings.TrimSpace(teachers[i].Fullname)
		if teachers[i].Fullname == "" {
			return models.ErrInvalidDataInput
		}

		normalized := strings.ToLower(teachers[i].Fullname)
		if _, ok := seen[normalized]; ok {
			return models.ErrAlreadyExists
		}
		seen[normalized] = struct{}{}

		exists, err := s.repo.TeacherExistsByFullname(teachers[i].Fullname)
		if err != nil {
			return err
		}
		if exists {
			return models.ErrAlreadyExists
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

	seen := make(map[string]struct{})

	for i := range classrooms {
		classrooms[i].Number = strings.TrimSpace(classrooms[i].Number)
		if classrooms[i].Number == "" {
			return models.ErrInvalidDataInput
		}

		normalized := strings.ToLower(classrooms[i].Number)
		if _, ok := seen[normalized]; ok {
			return models.ErrAlreadyExists
		}
		seen[normalized] = struct{}{}

		exists, err := s.repo.ClassroomExistsByNumber(classrooms[i].Number)
		if err != nil {
			return err
		}
		if exists {
			return models.ErrAlreadyExists
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

	seen := make(map[string]struct{})

	for i := range subjects {
		subjects[i].Name = strings.TrimSpace(subjects[i].Name)
		if subjects[i].Name == "" {
			return models.ErrInvalidDataInput
		}

		normalized := strings.ToLower(subjects[i].Name)
		if _, ok := seen[normalized]; ok {
			return models.ErrAlreadyExists
		}
		seen[normalized] = struct{}{}

		exists, err := s.repo.SubjectExistsByName(subjects[i].Name)
		if err != nil {
			return err
		}
		if exists {
			return models.ErrAlreadyExists
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

	seen := make(map[string]struct{})

	for i := range groups {
		groups[i].Name = strings.TrimSpace(groups[i].Name)
		if groups[i].Name == "" {
			return models.ErrInvalidDataInput
		}

		normalized := strings.ToLower(groups[i].Name)
		if _, ok := seen[normalized]; ok {
			return models.ErrAlreadyExists
		}
		seen[normalized] = struct{}{}

		exists, err := s.repo.GroupExistsByName(groups[i].Name)
		if err != nil {
			return err
		}
		if exists {
			return models.ErrAlreadyExists
		}
	}

	if err := s.repo.AddGroup(groups); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) GetTeachers() ([]models.Teacher, error) {
	return s.repo.GetTeachers()
}

func (s *AdminService) GetSubjects() ([]models.Subject, error) {
	return s.repo.GetSubjects()
}

func (s *AdminService) GetClassrooms() ([]models.Classroom, error) {
	return s.repo.GetClassrooms()
}

func (s *AdminService) GetGroups() ([]models.Group, error) {
	return s.repo.GetGroups()
}
