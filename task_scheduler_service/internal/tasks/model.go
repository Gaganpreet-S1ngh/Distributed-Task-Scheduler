package tasks

import (
	"time"

	"github.com/uptrace/bun"
)

type Task struct {
	bun.BaseModel `bun:"table:tasks,alias:t"`

	ID          uint64    `bun:"id,pk,autoincrement"`
	Command     string    `bun:"command,notnull"`
	ScheduledAt time.Time `bun:"scheduled_at,nullzero,notnull"`
	PickedAt    time.Time `bun:"picked_at"`
	StartedAt   time.Time `bun:"started_at"`
	CompletedAt time.Time `bun:"completed_at"`
	FailedAt    time.Time `bun:"failed_at"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}
