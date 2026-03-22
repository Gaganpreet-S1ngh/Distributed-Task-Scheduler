package tasks

import "net/http"

type Routes struct {
	handler *Handler
}

func NewRoutes(h *Handler) *Routes {
	return &Routes{handler: h}
}

func (r *Routes) SetupPublicRoutes() {
	http.HandleFunc("/schedule", r.handler.HandleScheduleTask)
	http.HandleFunc("/status", r.handler.HandleGetTaskStatus)
	http.HandleFunc("/health", HandleHealth)
}
