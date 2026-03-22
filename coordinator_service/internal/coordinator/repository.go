package coordinator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type Repository interface {
	UpdateTaskStatus(ctx context.Context, taskID uint64, timeStamp time.Time, column string) error
	FetchScheduledTasks(ctx context.Context, tx bun.Tx) ([]Task, error)
	MarkTaskAsPicked(ctx context.Context, tx bun.Tx, taskID uint64) error
	RunInTx(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error
}

type repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) Repository {
	return &repository{db: db}
}

func (r *repository) UpdateTaskStatus(ctx context.Context, taskID uint64, timeStamp time.Time, column string) error {
	query := fmt.Sprintf("%s = ?", column)
	_, err := r.db.NewUpdate().
		TableExpr("tasks").
		Set(query, timeStamp).
		Where("id = ?", taskID).Exec(ctx)

	return err
}

func (r *repository) FetchScheduledTasks(ctx context.Context, tx bun.Tx) ([]Task, error) {
	var tasks []Task

	err := tx.NewSelect().
		Model(&tasks).
		Where("scheduled_at < (NOW() + INTERVAL '30 seconds')").
		Where("(picked_at IS NULL OR picked_at = '0001-01-01 00:00:00+00')").
		OrderExpr("scheduled_at").
		For("UPDATE SKIP LOCKED").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("fetch tasks: %w", err)
	}

	return tasks, nil
}

func (r *repository) MarkTaskAsPicked(ctx context.Context, tx bun.Tx, taskID uint64) error {
	_, err := tx.NewUpdate().
		TableExpr("tasks").
		Set("picked_at = NOW()").
		Where("id = ?", taskID).
		Exec(ctx)
	return err
}

func (r *repository) RunInTx(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, fn)
}
