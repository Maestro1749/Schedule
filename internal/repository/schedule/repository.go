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
		GroupID int,
		WeekType int,
		Weekday int,
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
		WHERE s.group_id = $1 AND s.weekday = $2 AND s.week_type = $3 AND (s.subgroup IS NULL OR s.subgroup = $4)
		ORDER BY s.lesson_number;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, GroupID, Weekday, WeekType, Subgroup)
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
