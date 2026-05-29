package task

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetTaskById(ctx context.Context, id uuid.UUID) (*Task, error)
	GetTasksAssignedToUser(ctx context.Context, userId uuid.UUID) ([]*Task, error)
	GetTasksCreatedByUser(ctx context.Context, userId uuid.UUID) ([]*Task, error)
	AddTask(ctx context.Context, task *Task) error
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
}
