package tasks

import (
	"context"

	"github.com/uptrace/bun"
)

type Repository interface {
	CreateTask(ctx context.Context, task *Task) error
	FindTaskByID(ctx context.Context, id int64) (*Task, error)
}

type repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateTask(ctx context.Context, task *Task) error {
	_, err := r.db.NewInsert().
		Model(task).
		Exec(ctx)
	return err
}

func (r *repository) FindTaskByID(ctx context.Context, id int64) (*Task, error) {
	task := new(Task)

	err := r.db.NewSelect().
		Model(task).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return task, nil
}
