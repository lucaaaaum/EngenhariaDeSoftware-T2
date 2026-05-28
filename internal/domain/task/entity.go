package task

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID
	Title       string
	Description string
	AssignedTo  uuid.UUID
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	Status      TaskStatus
}

func NewTask(title, description string, createdBy uuid.UUID) (*Task, error) {
	if len(title) == 0 {
		return nil, errors.New("Task title cannot be empty")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, errors.Join(errors.New("Failed to generate task ID"), err)
	}

	return &Task{
		Id:          id,
		Title:       title,
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		Status:      Pending,
	}, nil
}
