package task

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID    `db:"id"`
	Title       string       `db:"title"`
	Description string       `db:"description"`
	AssignedTo  uuid.UUID    `db:"assigned_to"`
	CreatedBy   uuid.UUID    `db:"created_by"`
	CreatedAt   time.Time    `db:"created_at"`
	DueDate     *time.Time   `db:"due_date"`
	Status      TaskStatus   `db:"status"`
	Priority    TaskPriority `db:"priority"`
}

func NewTask(title, description string, createdBy uuid.UUID, priority TaskPriority, dueDate *time.Time) (*Task, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generating task id: %w", err)
	}

	return &Task{
		Id:          id,
		Title:       title,
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		DueDate:     dueDate,
		Status:      Pending,
		Priority:    priority,
	}, nil
}
