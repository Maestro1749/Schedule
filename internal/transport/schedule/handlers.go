package schedule_handler

import (
	"encoding/json"
	"net/http"
	"schedule/internal/models"
	schedule_service "schedule/internal/service/schedule"
	"strconv"

	"go.uber.org/zap"
)

type ScheduleHandler struct {
	service *schedule_service.ScheduleService
	logger  *zap.Logger
}

func NewScheduleHandler(service *schedule_service.ScheduleService, logger *zap.Logger) *ScheduleHandler {
	return &ScheduleHandler{service: service, logger: logger}
}

func (h *ScheduleHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	groupID, err := strconv.Atoi(query.Get("group_id"))
	if err != nil {
		http.Error(w, "invalid group_id", http.StatusBadRequest)
		return
	}

	weekType, err := strconv.Atoi(query.Get("week_type"))
	if err != nil {
		http.Error(w, "invalid week_type", http.StatusBadRequest)
		return
	}

	weekday, err := strconv.Atoi(query.Get("weekday"))
	if err != nil {
		http.Error(w, "invalid weekday", http.StatusBadRequest)
		return
	}

	var subgroup *int
	if subgroupStr := query.Get("subgroup"); subgroupStr != "" {
		sg, err := strconv.Atoi(subgroupStr)
		if err != nil {
			http.Error(w, "invalid subgroup", http.StatusBadRequest)
			return
		}
		subgroup = &sg
	}

	req := models.GetScheduleRequest{
		GroupID:  groupID,
		WeekType: weekType,
		Weekday:  weekday,
		Subgroup: subgroup,
	}

	schedule, err := h.service.GetSchedule(req)
	if err != nil {
		http.Error(w, "failed to get schedule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(schedule); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
