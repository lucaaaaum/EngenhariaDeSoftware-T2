package task

import "github.com/google/uuid"

type Repository interface {
	GetTaskById(id uuid.UUID) (*Task, error)
	GetTasksAssignedToUser(userId uuid.UUID) ([]*Task, error)
	GetTasksCreatedByUser(userId uuid.UUID) ([]*Task, error)
	AddTask(task *Task) error
	UpdateTask(task *Task) error
	DeleteTask(id uuid.UUID) error
}
