package task

import (
	"time"
	"tarefas/internal/domain/task"

	"github.com/google/uuid"
)

// CreateTaskCommand são os dados para criar uma tarefa.
// Recebidos via JSON no POST /tasks.
type CreateTaskCommand struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	CreatedBy   uuid.UUID        `json:"created_by"`
	Priority    task.TaskPriority `json:"priority"`    // 0=baixa, 1=média, 2=alta
	DueDate     *time.Time       `json:"due_date"`    // opcional — pode ser nulo
}

// UpdateTaskCommand são os dados para atualizar uma tarefa.
// Recebidos via JSON no PUT /tasks/{id}.
type UpdateTaskCommand struct {
	Id          uuid.UUID        `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	AssignedTo  uuid.UUID        `json:"assigned_to"`
	Status      task.TaskStatus  `json:"status"`
	Priority    task.TaskPriority `json:"priority"`
	DueDate     *time.Time       `json:"due_date"`
}

// TaskDto é o que retornamos ao cliente.
type TaskDto struct {
	Id          uuid.UUID        `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	AssignedTo  uuid.UUID        `json:"assigned_to"`
	CreatedBy   uuid.UUID        `json:"created_by"`
	CreatedAt   time.Time        `json:"created_at"`
	DueDate     *time.Time       `json:"due_date"`
	Status      task.TaskStatus  `json:"status"`
	Priority    task.TaskPriority `json:"priority"`
}

func NewTaskDto(t *task.Task) *TaskDto {
	return &TaskDto{
		Id:          t.Id,
		Title:       t.Title,
		Description: t.Description,
		AssignedTo:  t.AssignedTo,
		CreatedBy:   t.CreatedBy,
		CreatedAt:   t.CreatedAt,
		DueDate:     t.DueDate,
		Status:      t.Status,
		Priority:    t.Priority,
	}
}
