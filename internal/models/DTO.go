package models

type GetScheduleRequest struct {
	GroupID  int  `form:"group_id"`
	WeekType int  `form:"week_type"`
	Weekday  int  `form:"weekday"`
	Subgroup *int `form:"subgroup"`
}

type GetWeekScheduleRequest struct {
	GroupID  int  `form:"group_id"`
	WeekType *int `form:"week_type"`
	Subgroup *int `form:"subgroup"`
}

type ScheduleItemResponse struct {
	ID           int    `json:"id"`
	GroupID      int    `json:"group_id"`
	SubjectID    int    `json:"subject_id"`
	TeacherID    int    `json:"teacher_id"`
	ClassroomID  int    `json:"classroom_id"`
	Weekday      int    `json:"weekday"`
	LessonNumber int    `json:"lesson_number"`
	WeekType     *int   `json:"week_type"`
	Subgroup     *int   `json:"subgroup"`
	SubjectName  string `json:"subject_name"`
	TeacherName  string `json:"teacher_name"`
	ClassroomNum string `json:"classroom_num"`
	GroupName    string `json:"group_name"`
}

type GetScheduleResponse struct {
	Items []ScheduleItemResponse `json:"items"`
}

// Teachers
type AddTeacherDTO struct {
	Fullname string `json:"fullname"`
}

type CreateTeachersRequest []AddTeacherDTO

// Subjects
type AddSubjectDTO struct {
	Name string
}

type CreateSubjectRequest []AddSubjectDTO

// Classrooms
type AddClassroomDTO struct {
	Number string
}

type CreateClassroomRequest []AddClassroomDTO

// Groups
type AddGroupDTO struct {
	Name string
}

type CreateGroupRequest []AddGroupDTO

type CreateScheduleDTO struct {
	GroupID      int  `json:"group_id"`
	SubjectID    int  `json:"subject_id"`
	TeacherID    int  `json:"teacher_id"`
	ClassroomID  int  `json:"classroom_id"`
	Weekday      int  `json:"weekday"`
	LessonNumber int  `json:"lesson_number"`
	WeekType     *int `json:"week_type,omitempty"`
	Subgroup     *int `json:"subgroup,omitempty"`
}
