package schedule_repository

import (
	"context"
	"database/sql"
	"schedule/internal/models"
	"time"

	"go.uber.org/zap"
)

type ScheduleRepository interface {
	GetSchedule(
		ctx context.Context,
		GroupID int,
		WeekType int,
		Weekday int,
		Subgroup *int,
	) ([]models.ScheduleItemResponse, error)

	GetWeekSchedule(
		ctx context.Context,
		GroupID int,
		WeekType *int,
		Subgroup *int,
	) ([]models.ScheduleItemResponse, error)
}

type scheduleRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewScheduleRepo(db *sql.DB, logger *zap.Logger) ScheduleRepository {
	return &scheduleRepo{db: db, logger: logger}
}

func (r *scheduleRepo) GetSchedule(
	ctx context.Context,
	GroupID int,
	WeekType int,
	Weekday int,
	Subgroup *int,
) ([]models.ScheduleItemResponse, error) {
	query := `
		SELECT
			s.id,
			s.group_id,
			s.subject_id,
			s.teacher_id,
			s.classroom_id,
			s.weekday,
			s.lesson_number,
			s.week_type,
			s.subgroup,
			sub.name,
			t.fullname,
			c.num,
			g.name
		FROM schedule s
		JOIN Subjects sub ON s.subject_id = sub.id
		JOIN Teachers t ON s.teacher_id = t.id
		JOIN Classrooms c ON s.classroom_id = c.id
		JOIN Groups g ON s.group_id = g.id
		WHERE s.group_id = $1 AND s.weekday = $2 AND (s.week_type IS NULL OR s.week_type = $3) AND (s.subgroup IS NULL OR s.subgroup = $4)
		ORDER BY s.lesson_number;
	`

	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, query, GroupID, Weekday, WeekType, Subgroup)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var schedule []models.ScheduleItemResponse
	for rows.Next() {
		var s models.ScheduleItemResponse
		if err := rows.Scan(
			&s.ID,
			&s.GroupID,
			&s.SubjectID,
			&s.TeacherID,
			&s.ClassroomID,
			&s.Weekday,
			&s.LessonNumber,
			&s.WeekType,
			&s.Subgroup,
			&s.SubjectName,
			&s.TeacherName,
			&s.ClassroomNum,
			&s.GroupName,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err))
			return nil, err
		}
		schedule = append(schedule, s)
	}

	return schedule, nil
}

func (r *scheduleRepo) GetWeekSchedule(
	ctx context.Context,
	GroupID int,
	WeekType *int,
	Subgroup *int,
) ([]models.ScheduleItemResponse, error) {
	query := `
		SELECT
			s.id,
			s.group_id,
			s.subject_id,
			s.teacher_id,
			s.classroom_id,
			s.weekday,
			s.lesson_number,
			s.week_type,
			s.subgroup,
			sub.name,
			t.fullname,
			c.num,
			g.name
		FROM schedule s
		JOIN Subjects sub ON s.subject_id = sub.id
		JOIN Teachers t ON s.teacher_id = t.id
		JOIN Classrooms c ON s.classroom_id = c.id
		JOIN Groups g ON s.group_id = g.id
		WHERE s.group_id = $1
		  AND ($2::int IS NULL OR s.week_type IS NULL OR s.week_type = $2)
		  AND ($3::int IS NULL OR s.subgroup IS NULL OR s.subgroup = $3)
		ORDER BY s.weekday, s.lesson_number, s.week_type;
	`

	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, query, GroupID, WeekType, Subgroup)
	if err != nil {
		r.logger.Error("failed to execute weekly schedule query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var schedule []models.ScheduleItemResponse
	for rows.Next() {
		var s models.ScheduleItemResponse
		if err := rows.Scan(
			&s.ID,
			&s.GroupID,
			&s.SubjectID,
			&s.TeacherID,
			&s.ClassroomID,
			&s.Weekday,
			&s.LessonNumber,
			&s.WeekType,
			&s.Subgroup,
			&s.SubjectName,
			&s.TeacherName,
			&s.ClassroomNum,
			&s.GroupName,
		); err != nil {
			r.logger.Error("failed to scan weekly schedule row", zap.Error(err))
			return nil, err
		}
		schedule = append(schedule, s)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("weekly schedule rows iteration error", zap.Error(err))
		return nil, err
	}

	return schedule, nil
}
