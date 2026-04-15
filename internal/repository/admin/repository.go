package admin_repository

import (
	"context"
	"database/sql"
	"fmt"
	"schedule/internal/models"
	"time"

	"go.uber.org/zap"
)

type AdminRepository interface {
	AddTeacher(ctx context.Context, teachers []models.Teacher) error
	AddSubject(ctx context.Context, subjects []models.Subject) error
	AddClassroom(ctx context.Context, classrooms []models.Classroom) error
	AddGroup(ctx context.Context, groups []models.Group) error
	TeacherExistsByFullname(ctx context.Context, fullname string) (bool, error)
	SubjectExistsByName(ctx context.Context, name string) (bool, error)
	ClassroomExistsByNumber(ctx context.Context, number string) (bool, error)
	GroupExistsByName(ctx context.Context, name string) (bool, error)
	GetTeachers(ctx context.Context) ([]models.Teacher, error)
	GetSubjects(ctx context.Context) ([]models.Subject, error)
	GetClassrooms(ctx context.Context) ([]models.Classroom, error)
	GetGroups(ctx context.Context) ([]models.Group, error)
	CreateSchedule(ctx context.Context, data []models.CreateScheduleDTO) error
	DeleteSchedule(ctx context.Context, groupID int, weekday int, weektype *int, subgroup *int, lessonNumber *int) error
	GetGroupIdByName(ctx context.Context, groupName string) (int, error)
}

type adminRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewAdminRepository(db *sql.DB, logger *zap.Logger) AdminRepository {
	return &adminRepo{db: db, logger: logger}
}

func (r *adminRepo) CreateSchedule(ctx context.Context, data []models.CreateScheduleDTO) error {
	txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(txCtx, nil)
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(txCtx, "INSERT INTO Schedule (group_id, subject_id, teacher_id, classroom_id, weekday, lesson_number, week_type, subgroup) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")
	if err != nil {
		r.logger.Error("Failed to prepare statement", zap.Error(err))
		return models.ErrInternalServer
	}

	for _, item := range data {
		var weekType interface{}
		if item.WeekType != nil {
			weekType = *item.WeekType
		}

		if _, err := stmt.ExecContext(txCtx, item.GroupID, item.SubjectID, item.TeacherID, item.ClassroomID, item.Weekday, item.LessonNumber, weekType, item.Subgroup); err != nil {
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

func (r *adminRepo) AddTeacher(ctx context.Context, teachers []models.Teacher) error {
	txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(txCtx, nil)
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(txCtx, "INSERT INTO Teachers (fullname) VALUES ($1)")
	if err != nil {
		r.logger.Error("Failed to prepare statement", zap.Error(err))
		return models.ErrInternalServer
	}

	for _, t := range teachers {
		if _, err := stmt.ExecContext(txCtx, t.Fullname); err != nil {
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

func (r *adminRepo) TeacherExistsByFullname(ctx context.Context, fullname string) (bool, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool
	err := r.db.QueryRowContext(queryCtx, "SELECT EXISTS(SELECT 1 FROM Teachers WHERE LOWER(fullname) = LOWER($1))", fullname).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check teacher existence", zap.Error(err))
		return false, models.ErrInternalServer
	}

	return exists, nil
}

func (r *adminRepo) GetTeachers(ctx context.Context) ([]models.Teacher, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT id, fullname FROM Teachers ORDER BY fullname ASC")
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

func (r *adminRepo) AddSubject(ctx context.Context, subjects []models.Subject) error {
	txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(txCtx, nil)
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(txCtx, "INSERT INTO Subjects (name) VALUES ($1)")
	if err != nil {
		r.logger.Error("Failed to prepare statement", zap.Error(err))
		return models.ErrInternalServer
	}

	for _, s := range subjects {
		if _, err := stmt.ExecContext(txCtx, s.Name); err != nil {
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

func (r *adminRepo) SubjectExistsByName(ctx context.Context, name string) (bool, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool
	err := r.db.QueryRowContext(queryCtx, "SELECT EXISTS(SELECT 1 FROM Subjects WHERE LOWER(name) = LOWER($1))", name).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check subject existence", zap.Error(err))
		return false, models.ErrInternalServer
	}

	return exists, nil
}

func (r *adminRepo) GetSubjects(ctx context.Context) ([]models.Subject, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT id, name FROM Subjects ORDER BY name ASC")
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

func (r *adminRepo) AddClassroom(ctx context.Context, classrooms []models.Classroom) error {
	txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(txCtx, nil)
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(txCtx, "INSERT INTO Classrooms (num) VALUES ($1)")
	if err != nil {
		r.logger.Error("Failed to prepare statement", zap.Error(err))
		return models.ErrInternalServer
	}

	for _, c := range classrooms {
		if _, err := stmt.ExecContext(txCtx, c.Number); err != nil {
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

func (r *adminRepo) ClassroomExistsByNumber(ctx context.Context, number string) (bool, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool
	err := r.db.QueryRowContext(queryCtx, "SELECT EXISTS(SELECT 1 FROM Classrooms WHERE LOWER(num) = LOWER($1))", number).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check classroom existence", zap.Error(err))
		return false, models.ErrInternalServer
	}

	return exists, nil
}

func (r *adminRepo) GetClassrooms(ctx context.Context) ([]models.Classroom, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT id, num FROM Classrooms ORDER BY num ASC")
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

func (r *adminRepo) AddGroup(ctx context.Context, groups []models.Group) error {
	txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(txCtx, nil)
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return models.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(txCtx, "INSERT INTO Groups (name) VALUES ($1)")
	if err != nil {
		r.logger.Error("Failed to prepare statement", zap.Error(err))
		return models.ErrInternalServer
	}

	for _, g := range groups {
		if _, err := stmt.ExecContext(txCtx, g.Name); err != nil {
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

func (r *adminRepo) GroupExistsByName(ctx context.Context, name string) (bool, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool
	err := r.db.QueryRowContext(queryCtx, "SELECT EXISTS(SELECT 1 FROM Groups WHERE LOWER(name) = LOWER($1))", name).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check group existence", zap.Error(err))
		return false, models.ErrInternalServer
	}

	return exists, nil
}

func (r *adminRepo) GetGroups(ctx context.Context) ([]models.Group, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT id, name FROM Groups ORDER BY name ASC")
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

func (r *adminRepo) DeleteSchedule(
	ctx context.Context,
	groupID int,
	weekday int,
	weektype *int,
	subgroup *int,
	lessonNumber *int,
) error {
	query := "DELETE FROM Schedule WHERE group_id = $1 AND weekday = $2"
	args := []interface{}{groupID, weekday}
	argIndex := 3

	if weektype != nil {
		query += fmt.Sprintf(" AND week_type = $%d", argIndex)
		args = append(args, *weektype)
		argIndex++
	}

	if subgroup != nil {
		query += fmt.Sprintf(" AND subgroup = $%d", argIndex)
		args = append(args, *subgroup)
		argIndex++
	}

	if lessonNumber != nil {
		query += fmt.Sprintf(" AND lesson_number = $%d", argIndex)
		args = append(args, *lessonNumber)
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, query, args...)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err))
		return models.ErrInternalServer
	}

	count, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("error to get affected rows", zap.Error(err))
		return models.ErrInternalServer
	}
	if count == 0 {
		return models.ErrNotUpdated
	}

	return nil
}

func (r *adminRepo) GetGroupIdByName(ctx context.Context, groupName string) (int, error) {
	query := "SELECT id FROM Groups WHERE name = $1"

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var id int

	err := r.db.QueryRowContext(queryCtx, query, groupName).Scan(&id)
	if err != nil {
		r.logger.Error("error to execute query", zap.Error(err))
		return 0, models.ErrInternalServer
	}

	return id, nil
}
