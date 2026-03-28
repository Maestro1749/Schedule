package admin_repository

import (
	"database/sql"
	"schedule/internal/models"

	"go.uber.org/zap"
)

type AdminRepository interface {
	AddTeacher(teachers []models.Teacher) error
	AddSubject(subjects []models.Subject) error
	AddClassroom(classrooms []models.Classroom) error
	AddGroup(groups []models.Group) error
	CreateSchedule(data []models.CreateScheduleDTO) error
}

type adminRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewAdminRepository(db *sql.DB, logger *zap.Logger) AdminRepository {
	return &adminRepo{db: db, logger: logger}
}

func (r *adminRepo) CreateSchedule(data []models.CreateScheduleDTO) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO Schedule (group_id, subject_id, teacher_id, classroom_id, weekday, lesson_number, week_type, subgroup) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")
	if err != nil {
		r.logger.Error("Failed to prepare statement", zap.Error(err))
		return models.ErrInternalServer
	}

	for _, item := range data {
		if _, err := stmt.Exec(item.GroupID, item.SubjectID, item.TeacherID, item.ClassroomID, item.Weekday, item.LessonNumber, item.WeekType, item.Subgroup); err != nil {
			r.logger.Error("Failed to execute statement", zap.Error(err))
			return models.ErrInternalServer
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Failed to commit transaction", zap.Error(err))
		return models.ErrInternalServer
	}

	r.logger.Info("Successfully created schedule", zap.Int("count", len(data)))
	return nil
}

func (r *adminRepo) AddTeacher(teachers []models.Teacher) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO Teachers (fullname) VALUES ($1)")
	if err != nil {
		r.logger.Error("Failed to prepare statement", zap.Error(err))
		return models.ErrInternalServer
	}

	for _, t := range teachers {
		if _, err := stmt.Exec(t.Fullname); err != nil {
			r.logger.Error("Failed to execute statement", zap.Error(err))
			return models.ErrInternalServer
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Failed to commit transaction", zap.Error(err))
		return models.ErrInternalServer
	}

	r.logger.Info("Successfully added teachers", zap.Int("count", len(teachers)))
	return nil
}

func (r *adminRepo) AddSubject(subjects []models.Subject) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO Subjects (name) VALUES ($1)")
	if err != nil {
		r.logger.Error("Failed to prepare statement", zap.Error(err))
		return models.ErrInternalServer
	}

	for _, s := range subjects {
		if _, err := stmt.Exec(s.Name); err != nil {
			r.logger.Error("Failed to execute statement", zap.Error(err))
			return models.ErrInternalServer
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Failed to commit transaction", zap.Error(err))
		return models.ErrInternalServer
	}

	r.logger.Info("Successfully added subjects", zap.Int("count", len(subjects)))
	return nil
}

func (r *adminRepo) AddClassroom(classrooms []models.Classroom) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO Classrooms (num) VALUES ($1)")
	if err != nil {
		r.logger.Error("Failed to prepare statement", zap.Error(err))
		return models.ErrInternalServer
	}

	for _, c := range classrooms {
		if _, err := stmt.Exec(c.Number); err != nil {
			r.logger.Error("Failed to execute statement", zap.Error(err))
			return models.ErrInternalServer
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Failed to commit transaction", zap.Error(err))
		return models.ErrInternalServer
	}

	r.logger.Info("Successfully added classrooms", zap.Int("count", len(classrooms)))
	return nil
}

func (r *adminRepo) AddGroup(groups []models.Group) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO Groups (name) VALUES ($1)")
	if err != nil {
		r.logger.Error("Failed to prepare statement", zap.Error(err))
		return models.ErrInternalServer
	}

	for _, g := range groups {
		if _, err := stmt.Exec(g.Name); err != nil {
			r.logger.Error("Failed to execute statement", zap.Error(err))
			return models.ErrInternalServer
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Failed to commit transaction", zap.Error(err))
		return models.ErrInternalServer
	}

	r.logger.Info("Successfully added groups", zap.Int("count", len(groups)))
	return nil
}
