package models

type Schedule struct {
	ID           int
	GroupID      int
	SubjectID    int
	TeacherID    int
	ClassroomID  int
	Weekday      int
	LessonNumber int
	WeekType     int
	Subgroup     *int
}

type Teacher struct {
	ID       int
	Fullname string
}

type Subject struct {
	ID   int
	Name string
}

type Classroom struct {
	ID     int
	Number string
}

type Group struct {
	ID   int
	Name string
}
