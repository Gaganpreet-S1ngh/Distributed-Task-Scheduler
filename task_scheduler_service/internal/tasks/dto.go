package tasks

type TaskResponse struct {
	TaskID      int64  `json:"task_id"`
	Command     string `json:"command"`
	ScheduledAt string `json:"scheduled_at,omitempty"`
	PickedAt    string `json:"picked_at,omitempty"`
	StartedAt   string `json:"started_at,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
	FailedAt    string `json:"failed_at,omitempty"`
}

type CreateTaskRequest struct {
	Command     string `json:"command" validate:"required"`
	ScheduledAt string `json:"scheduled_at" validate:"required"`
}

type CreateTaskResponse struct {
	TaskID      uint64 `json:"task_id"`
	Command     string `json:"command"`
	ScheduledAt string `json:"scheduled_at"`
}
