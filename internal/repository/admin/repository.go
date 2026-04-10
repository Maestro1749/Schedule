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
	TeacherExistsByFullname(fullname string) (bool, error)
	SubjectExistsByName(name string) (bool, error)
	ClassroomExistsByNumber(number string) (bool, error)
	GroupExistsByName(name string) (bool, error)
	GetTeachers() ([]models.Teacher, error)
	GetSubjects() ([]models.Subject, error)
	GetClassrooms() ([]models.Classroom, error)
	GetGroups() ([]models.Group, error)
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
		var weekType interface{}
		if item.WeekType != nil {
			weekType = *item.WeekType
		}

		if _, err := stmt.Exec(item.GroupID, item.SubjectID, item.TeacherID, item.ClassroomID, item.Weekday, item.LessonNumber, weekType, item.Subgroup); err != nil {
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

func (r *adminRepo) TeacherExistsByFullname(fullname string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM Teachers WHERE LOWER(fullname) = LOWER($1))", fullname).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check teacher existence", zap.Error(err))
		return false, models.ErrInternalServer
	}

	return exists, nil
}

func (r *adminRepo) GetTeachers() ([]models.Teacher, error) {
	rows, err := r.db.Query("SELECT id, fullname FROM Teachers ORDER BY fullname ASC")
	if err != nil {
		r.logger.Error("Failed to query teachers", zap.Error(err))
		return nil, models.ErrInternalServer
	}
	defer rows.Close()

	teachers := make([]models.Teacher, 0)
	for rows.Next() {
		var teacher models.Teacher
		if err := rows.Scan(&teacher.ID, &teacher.Fullname); err != nil {
			r.logger.Error("Failed to scan teacher row", zap.Error(err))
			return nil, models.ErrInternalServer
		}
		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Rows iteration error for teachers", zap.Error(err))
		return nil, models.ErrInternalServer
	}

	return teachers, nil
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

func (r *adminRepo) SubjectExistsByName(name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM Subjects WHERE LOWER(name) = LOWER($1))", name).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check subject existence", zap.Error(err))
		return false, models.ErrInternalServer
	}

	return exists, nil
}

func (r *adminRepo) GetSubjects() ([]models.Subject, error) {
	rows, err := r.db.Query("SELECT id, name FROM Subjects ORDER BY name ASC")
	if err != nil {
		r.logger.Error("Failed to query subjects", zap.Error(err))
		return nil, models.ErrInternalServer
	}
	defer rows.Close()

	subjects := make([]models.Subject, 0)
	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.ID, &subject.Name); err != nil {
			r.logger.Error("Failed to scan subject row", zap.Error(err))
			return nil, models.ErrInternalServer
		}
		subjects = append(subjects, subject)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Rows iteration error for subjects", zap.Error(err))
		return nil, models.ErrInternalServer
	}

	return subjects, nil
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

func (r *adminRepo) ClassroomExistsByNumber(number string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM Classrooms WHERE LOWER(num) = LOWER($1))", number).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check classroom existence", zap.Error(err))
		return false, models.ErrInternalServer
	}

	return exists, nil
}

func (r *adminRepo) GetClassrooms() ([]models.Classroom, error) {
	rows, err := r.db.Query("SELECT id, num FROM Classrooms ORDER BY num ASC")
	if err != nil {
		r.logger.Error("Failed to query classrooms", zap.Error(err))
		return nil, models.ErrInternalServer
	}
	defer rows.Close()

	classrooms := make([]models.Classroom, 0)
	for rows.Next() {
		var classroom models.Classroom
		if err := rows.Scan(&classroom.ID, &classroom.Number); err != nil {
			r.logger.Error("Failed to scan classroom row", zap.Error(err))
			return nil, models.ErrInternalServer
		}
		classrooms = append(classrooms, classroom)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Rows iteration error for classrooms", zap.Error(err))
		return nil, models.ErrInternalServer
	}

	return classrooms, nil
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

func (r *adminRepo) GroupExistsByName(name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM Groups WHERE LOWER(name) = LOWER($1))", name).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check group existence", zap.Error(err))
		return false, models.ErrInternalServer
	}

	return exists, nil
}

func (r *adminRepo) GetGroups() ([]models.Group, error) {
	rows, err := r.db.Query("SELECT id, name FROM Groups ORDER BY name ASC")
	if err != nil {
		r.logger.Error("Failed to query groups", zap.Error(err))
		return nil, models.ErrInternalServer
	}
	defer rows.Close()

	groups := make([]models.Group, 0)
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			r.logger.Error("Failed to scan group row", zap.Error(err))
			return nil, models.ErrInternalServer
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Rows iteration error for groups", zap.Error(err))
		return nil, models.ErrInternalServer
	}

	return groups, nil
}
