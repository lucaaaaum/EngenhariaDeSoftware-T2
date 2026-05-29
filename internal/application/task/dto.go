package task

import (
	"tarefas/internal/domain/task"

	"github.com/google/uuid"
)

type CreateTaskCommand struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"created_by"`
}

type UpdateTaskCommand struct {
	Id          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	AssignedTo  uuid.UUID       `json:"assigned_to"`
	Status      task.TaskStatus `json:"status"`
}

type TaskDto struct {
	Id          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	AssignedTo  uuid.UUID       `json:"assigned_to"`
	CreatedBy   uuid.UUID       `json:"created_by"`
	Status      task.TaskStatus `json:"status"`
}

func NewUserDto(task *task.Task) *TaskDto {
	return &TaskDto{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		AssignedTo:  task.AssignedTo,
		CreatedBy:   task.CreatedBy,
		Status:      task.Status,
	}
}
