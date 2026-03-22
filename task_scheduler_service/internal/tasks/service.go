package tasks

import "context"

type Service interface {
	CreateNewTask(ctx context.Context, t *Task) (uint64, error)
	GetTaskStatus(ctx context.Context, taskID int64) (Task, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateNewTask(ctx context.Context, t *Task) (uint64, error) {
	err := s.repo.CreateTask(ctx, t)
	return t.ID, err
}

func (s *service) GetTaskStatus(ctx context.Context, taskID int64) (Task, error) {
	result, err := s.repo.FindTaskByID(ctx, taskID)
	if err != nil {
		return Task{}, err
	}
	return *result, nil
}
