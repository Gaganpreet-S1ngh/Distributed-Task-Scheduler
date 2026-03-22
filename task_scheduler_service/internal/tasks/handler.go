package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

/* HELPER FUNCTIONS */

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

/* HANDLER FUNCTIONS (HTTP) */

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "only GET method is allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "UP",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Handler) HandleScheduleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeError(w, http.StatusMethodNotAllowed, "only POST method is allowed")
		return
	}

	var req CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Remove in production never do printing in Application layer because i/o
	log.Printf("Received schedule request: %+v", req)

	// 2025-01-01T00:00:00Z -------> 2025-01-01 00:00:00 +0000 UTC
	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledAt)

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format, use RFC3339 e.g. 2025-01-01T00:00:00Z")
		return
	}

	// Calling Service

	task := &Task{
		Command:     req.Command,
		ScheduledAt: scheduledTime,
	}

	taskID, err := s.svc.CreateNewTask(r.Context(), task)

	if err != nil {
		log.Printf("CreateNewTask error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to submit task")
		return
	}

	response := CreateTaskResponse{
		Command:     req.Command,
		ScheduledAt: req.ScheduledAt,
		TaskID:      taskID,
	}

	writeJSON(w, http.StatusCreated, response)

}

func (s *Handler) HandleGetTaskStatus(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		writeError(w, http.StatusMethodNotAllowed, "only GET method is allowed")
		return
	}

	taskIDStr := r.URL.Query().Get("task_id")

	if taskIDStr == "" {
		writeError(w, http.StatusBadRequest, "task_id is required")
		return
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task_id, must be a number")
		return
	}

	task, err := s.svc.GetTaskStatus(r.Context(), taskID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get task status. Error : %s", err), http.StatusInternalServerError)
		return
	}

	response := TaskResponse{
		TaskID:      taskID,
		Command:     task.Command,
		ScheduledAt: task.ScheduledAt.String(),
		PickedAt:    task.PickedAt.String(),
		StartedAt:   task.StartedAt.String(),
		CompletedAt: task.CompletedAt.String(),
		FailedAt:    task.FailedAt.String(),
	}

	writeJSON(w, http.StatusOK, response)
}
