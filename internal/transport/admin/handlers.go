package admin_handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"schedule/internal/models"
	admin_service "schedule/internal/service/admin"

	"go.uber.org/zap"
)

type AdminHandler struct {
	service *admin_service.AdminService
	logger  *zap.Logger
}

func NewUserHandler(service *admin_service.AdminService, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{service: service, logger: logger}
}

func (h *AdminHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var req []models.CreateScheduleDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateSchedule(req); err != nil {
		http.Error(w, "Failed to create schedule", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AdminHandler) AddTeacher(w http.ResponseWriter, r *http.Request) {
	var data models.CreateTeachersRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var teachers []models.Teacher
	for _, t := range data {
		teachers = append(teachers, models.Teacher{Fullname: t.Fullname})
	}

	if err := h.service.AddTeacher(teachers); err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidDataInput):
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AdminHandler) AddSubject(w http.ResponseWriter, r *http.Request) {
	var data models.CreateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var subjects []models.Subject
	for _, s := range data {
		subjects = append(subjects, models.Subject{Name: s.Name})
	}

	if err := h.service.AddSubject(subjects); err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidDataInput):
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AdminHandler) AddClassroom(w http.ResponseWriter, r *http.Request) {
	var data models.CreateClassroomRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var classrooms []models.Classroom
	for _, s := range data {
		classrooms = append(classrooms, models.Classroom{Number: s.Number})
	}

	if err := h.service.AddClassroom(classrooms); err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidDataInput):
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AdminHandler) AddGroup(w http.ResponseWriter, r *http.Request) {
	var data models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var groups []models.Group
	for _, s := range data {
		groups = append(groups, models.Group{Name: s.Name})
	}

	if err := h.service.AddGroup(groups); err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidDataInput):
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}
